package levelstore

import (
	"os"

	"github.com/nulayer/chainstore"
	"github.com/syndtr/goleveldb/leveldb"
)

type levelStore struct {
	storePath string
	db        *leveldb.DB
}

func New(storePath string) (*levelStore, error) {
	store := &levelStore{storePath: storePath}
	err := store.Open()
	return store, err
}

func (s *levelStore) Open() (err error) {
	// Create the store directory if doesnt exist
	if _, err = os.Stat(s.storePath); os.IsNotExist(err) {
		err = os.MkdirAll(s.storePath, 0755)
		if err != nil {
			return
		}
	}

	s.db, err = leveldb.OpenFile(s.storePath, nil)
	return
}

func (s *levelStore) Close() error {
	return s.db.Close()
}

func (s *levelStore) Put(key string, obj []byte) error {
	if !chainstore.IsValidKey(key) {
		return chainstore.ErrInvalidKey
	}
	return s.db.Put([]byte(key), obj, nil)
}

func (s *levelStore) Get(key string) ([]byte, error) {
	if !chainstore.IsValidKey(key) {
		return nil, chainstore.ErrInvalidKey
	}

	obj, err := s.db.Get([]byte(key), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return nil, err
	}
	return obj, nil
}

func (s *levelStore) Del(key string) error {
	if !chainstore.IsValidKey(key) {
		return chainstore.ErrInvalidKey
	}
	return s.db.Delete([]byte(key), nil)
}
