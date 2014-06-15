package memstore

import (
	"sync"

	"github.com/nulayer/chainstore"
	"github.com/nulayer/chainstore/lrumgr"
)

type memStore struct {
	sync.Mutex
	data map[string][]byte
}

// TODO: can we have our own middleware..? .. then it would be lruManager -> memCache

func New(capacity int64) *lrumgr.LruManager {
	memStore := &memStore{}
	memStore.Open() // TODO............
	store := lrumgr.New(capacity, memStore)
	return store
}

func (s *memStore) Open() (err error) {
	s.data = make(map[string][]byte, 1000)
	return nil
}

func (s *memStore) Close() error {
	return nil
}

func (s *memStore) Put(key string, obj []byte) (err error) {
	if !chainstore.IsValidKey(key) {
		return chainstore.ErrInvalidKey
	}
	s.Lock()
	s.data[key] = obj
	s.Unlock()
	return nil
}

func (s *memStore) Get(key string) (obj []byte, err error) {
	return s.data[key], nil
}

func (s *memStore) Del(key string) (err error) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()
	return nil
}
