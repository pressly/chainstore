package chainstore

import (
	"sync"
)

type memCache struct {
	sync.Mutex
	data map[string][]byte
}

// TODO: can we have our own middleware..? .. then it would be lruManager -> memCache

func MemCacheStore(capacity int64) (store *lruManager) {
	mem := &memCache{}
	mem.Open() // TODO............
	store = LRUManager(mem, capacity)
	return
}

func (s *memCache) Open() (err error) {
	s.data = make(map[string][]byte, 1000)
	return nil
}

func (s *memCache) Close() error {
	return nil
}

func (s *memCache) Put(key string, obj []byte) (err error) {
	if !IsValidKey(key) {
		return ErrInvalidKey
	}
	s.Lock()
	s.data[key] = obj
	s.Unlock()
	return nil
}

func (s *memCache) Get(key string) (obj []byte, err error) {
	return s.data[key], nil
}

func (s *memCache) Del(key string) (err error) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()
	return nil
}
