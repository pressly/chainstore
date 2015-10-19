package mockstore

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/pressly/chainstore"
	"golang.org/x/net/context"
)

var _ = chainstore.Store(&mockStore{})

type mockStore struct {
	mu      sync.RWMutex
	data    map[string][]byte
	cfg     *Config
	timer   *time.Timer
	timerMu sync.Mutex
	closed  bool
}

// Config stores settings for this store.
type Config struct {
	Capacity    int64
	SuccessRate float32
	Delay       time.Duration
}

// New creates and returns a mock chainstore.Store.
func New(cfg *Config) chainstore.Store {
	if cfg == nil {
		cfg = &Config{
			Capacity:    1000,
			SuccessRate: 1.0,
			Delay:       0,
		}
	}
	mockStore := &mockStore{
		data: make(map[string][]byte, cfg.Capacity),
		cfg:  cfg,
	}
	return mockStore
}

func (s *mockStore) success() bool {
	return rand.Float32() < s.cfg.SuccessRate
}

func (s *mockStore) Open() error {
	if !s.success() {
		return errors.New("Failed to open: random fail.")
	}
	return nil
}

func (s *mockStore) Close() error {
	if !s.success() {
		return errors.New("Failed to close: random fail.")
	}

	s.timerMu.Lock()
	s.timer.Reset(0)
	s.timerMu.Unlock()

	s.mu.Lock()
	s.closed = true
	s.data = nil
	s.mu.Unlock()

	return nil
}

func (s *mockStore) delay() {
	s.timerMu.Lock()
	s.timer = time.NewTimer(s.cfg.Delay)
	s.timerMu.Unlock()

	<-s.timer.C
}

func (s *mockStore) Put(ctx context.Context, key string, val []byte) (err error) {
	if !s.success() {
		return errors.New("Failed to put key in store: random fail.")
	}

	s.delay()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return errors.New("Store is closed.")
	}
	s.data[key] = val

	return nil
}

func (s *mockStore) Get(ctx context.Context, key string) ([]byte, error) {
	if !s.success() {
		return nil, errors.New("Failed to get key from store: random fail.")
	}

	s.delay()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, errors.New("Store is closed.")
	}
	val, ok := s.data[key]

	if !ok {
		return nil, fmt.Errorf("No such key %q", key)
	}

	return val, nil
}

func (s *mockStore) Del(ctx context.Context, key string) error {
	if !s.success() {
		return errors.New("Failed to delete key from store: random fail.")
	}

	s.delay()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return errors.New("Store is closed.")
	}

	if _, ok := s.data[key]; !ok {
		return fmt.Errorf("No such key %q", key)
	}

	delete(s.data, key)

	return nil
}
