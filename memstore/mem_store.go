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

func New(capacity int64) *lrumgr.LruManager {
	memStore := &memStore{}
	memStore.Open() // TODO: in case of error..?
	store := lrumgr.New(capacity, memStore)
	return store
}

func (s *memStore) Open() (err error) {
	s.data = make(map[string][]byte, 1000)
	return
}

func (s *memStore) Close() (err error) { return }

func (s *memStore) Put(key string, val []byte) (err error) {
	if !chainstore.IsValidKey(key) {
		return chainstore.ErrInvalidKey
	}
	s.Lock()
	s.data[key] = val
	s.Unlock()
	return nil
}

func (s *memStore) Get(key string) (val []byte, err error) {
	return s.data[key], nil
}

func (s *memStore) Del(key string) (err error) {
	s.Lock()
	delete(s.data, key)
	s.Unlock()
	return
}
