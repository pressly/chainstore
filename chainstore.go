package chainstore

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"golang.org/x/net/context"
)

type storeFn func(s Store) error

var (
	keyInvalidator = regexp.MustCompile(`(i?)[^a-z0-9\/_\-:\.]`)
)

var (
	DefaultTimeout = time.Millisecond * 3500
)

const (
	maxKeyLen = 256
)

type Store interface {
	Open() error
	Close() error
	Put(ctx context.Context, key string, val []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
}

type storeWrapper struct {
	Store
	errE  error
	errMu sync.RWMutex
}

func (s *storeWrapper) err() error {
	s.errMu.RLock()
	defer s.errMu.RUnlock()
	return s.errE
}

func (s *storeWrapper) setErr(err error) {
	s.errMu.Lock()
	defer s.errMu.Unlock()
	s.errE = err
}

// Chain represents a store chain.
type Chain struct {
	stores []*storeWrapper
}

// New creates a new store chain backed by the passed stores.
func New(stores ...Store) Store {
	c := &Chain{
		stores: make([]*storeWrapper, 0, len(stores)),
	}
	for _, s := range stores {
		c.stores = append(c.stores, &storeWrapper{Store: s})
	}
	return c
}

// Open opens all the stores.
func (c *Chain) Open() error {

	if err := c.firstErr(); err != nil {
		return fmt.Errorf("Open failed due to a previous error: %q", err)
	}

	var wg sync.WaitGroup

	for i := range c.stores {
		wg.Add(1)
		go func(s *storeWrapper) {
			defer wg.Done()
			s.setErr(s.Open())
		}(c.stores[i])
	}

	wg.Wait()

	return c.firstErr()
}

// Close closes all the stores.
func (c *Chain) Close() error {
	var wg sync.WaitGroup

	for i := range c.stores {
		wg.Add(1)
		go func(s *storeWrapper) {
			defer wg.Done()
			s.setErr(s.Close())
		}(c.stores[i])
	}

	wg.Wait()

	return c.firstErr()
}

func (c *Chain) Put(ctx context.Context, key string, val []byte) (err error) {
	if !isValidKey(key) {
		return ErrInvalidKey
	}

	if err := c.firstErr(); err != nil {
		return fmt.Errorf("Open failed due to a previous error: %q", err)
	}

	fn := func(s Store) error {
		return s.Put(ctx, key, val)
	}

	return c.doWithContext(ctx, fn)
}

func (c *Chain) Get(ctx context.Context, key string) (val []byte, err error) {
	if !isValidKey(key) {
		return nil, ErrInvalidKey
	}

	if err := c.firstErr(); err != nil {
		return nil, fmt.Errorf("Open failed due to a previous error: %q", err)
	}

	nextStore := make(chan Store, len(c.stores))
	for _, store := range c.stores {
		nextStore <- store
	}
	close(nextStore)

	putBack := make(chan Store, len(c.stores))

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case store, ok := <-nextStore:
			if !ok {
				return nil, ErrNoSuchKey
			}

			val, err := store.Get(ctx, key)

			if err != nil {
				if err == ErrTimeout {
					return nil, err
				}
				putBack <- store
				continue
			}

			close(putBack)

			for store := range putBack {
				go store.Put(ctx, key, val)
			}

			return val, nil
		}
	}

	panic("reached")
}

func (c *Chain) Del(ctx context.Context, key string) (err error) {
	if !isValidKey(key) {
		return ErrInvalidKey
	}

	if err := c.firstErr(); err != nil {
		return fmt.Errorf("Delete failed due to a previous error: %q", err)
	}

	fn := func(s Store) error {
		return s.Del(ctx, key)
	}

	return c.doWithContext(ctx, fn)
}

func (c *Chain) doWithContext(ctx context.Context, fn storeFn) error {
	var wg sync.WaitGroup

	for i := range c.stores {
		wg.Add(1)

		go func(s *storeWrapper) {
			defer wg.Done()
			s.setErr(fn(s))
		}(c.stores[i])
	}

	wg.Wait()

	return c.firstErr()
}

func (c *Chain) firstErr() error {
	var rerr error
	for i := range c.stores {
		if err := c.stores[i].err(); err != nil {
			rerr = err
			if err == ErrTimeout {
				// We can recover from this kind of error, so we return it and try
				// again.
				c.stores[i].setErr(nil)
				return err
			} else {
				break
			}
		}
	}
	return rerr
}

func isValidKey(key string) bool {
	return len(key) <= maxKeyLen && !keyInvalidator.MatchString(key)
}
