package chainstore

// TODO: ... move this all into chainstore.go

// friggen cleannnnnnnnnnnnnnnnnnnn :D

type Chain struct {
	stores []Store
}

func NewChain(stores ...Store) (Store, error) {
	c := &Chain{}
	if len(stores) == 0 {
		return nil, ErrNoStores
	}
	c.stores = stores
	err := c.Open()
	return c, err
}

func (s *Chain) Open() (err error) {
	for _, store := range s.stores {
		err := store.Open()
		if err != nil {
			return err // return first error that comes up
		}
	}
	return // noop, we expect all of the stores to be opened on their own
}

func (s *Chain) Close() error {
	for _, store := range s.stores {
		err := store.Close()
		if err != nil {
			return err // return first error that comes up
		}
	}
	return nil
}

// TODO: better error support if an async Put() fails in some store.
// A channel for clients to listen in on for errors?
func (s *Chain) Put(key string, obj []byte) error {
	putFn := func(store Store, key string, obj []byte) error {
		return store.Put(key, obj)
	}

	for i, store := range s.stores {
		if i == 0 {
			err := putFn(store, key, obj)
			if err != nil {
				return err
			}
		} else {
			// TODO: should we group all async Puts in a single goroutine?
			go putFn(store, key, obj) // all other stores are async
		}
	}
	return nil
}

func (s *Chain) Get(key string) ([]byte, error) {
	for i, store := range s.stores {
		obj, err := store.Get(key) // return the first one that matches
		if err != nil {
			return nil, err
		}

		if len(obj) > 0 {
			if i > 0 { // save the value in all other stores up the chain
				go func() {
					for n := i - 1; n >= 0; n-- {
						s.stores[n].Put(key, obj)
					}
				}()
			}

			// return the object from the store
			return obj, nil
		}
	}
	return nil, nil
}

func (s *Chain) Del(key string) error {
	for i, store := range s.stores {
		if i == 0 {
			err := store.Del(key)
			if err != nil {
				return err
			}
		} else {
			go func() {
				s.stores[i].Del(key)
			}()
		}

	}
	return nil
}

// TODO:
// func (s *Chain) SetLogger(l *log.Logger) {
// }
