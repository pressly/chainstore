package leveldb

import (
	"os"

	. "github.com/nulayer/chainstore"
	"github.com/syndtr/goleveldb/leveldb"
)

type Leveldb struct {
	storePath string
	db        *leveldb.DB
}

func NewStore(storePath string) (store *Leveldb, err error) {
	store = &Leveldb{storePath: storePath}
	err = store.Open()
	return
}

func (s *Leveldb) Open() (err error) {
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

func (s *Leveldb) Close() error {
	return s.db.Close()
}

func (s *Leveldb) Put(key string, obj []byte) error {
	if !IsValidKey(key) {
		return ErrInvalidKey
	}
	return s.db.Put([]byte(key), obj, nil)
}

func (s *Leveldb) Get(key string) ([]byte, error) {
	if !IsValidKey(key) {
		return nil, ErrInvalidKey
	}

	obj, err := s.db.Get([]byte(key), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return nil, err
	}
	return obj, nil
}

func (s *Leveldb) Del(key string) error {
	if !IsValidKey(key) {
		return ErrInvalidKey
	}
	return s.db.Delete([]byte(key), nil)
}
