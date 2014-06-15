package boltstore

import (
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/nulayer/chainstore"
	"github.com/rcrowley/go-metrics"
)

// TODO: remove the metrics out of here...

type boltStore struct {
	storePath  string
	bucketName []byte

	db     *bolt.DB
	bucket *bolt.Bucket
}

func New(storePath string, bucketName string) (*boltStore, error) {
	store := &boltStore{storePath: storePath, bucketName: []byte(bucketName)}
	err := store.Open()
	return store, err
}

func (s *boltStore) Open() (err error) {
	// Create the store directory if doesnt exist
	storeDir := filepath.Dir(s.storePath)
	if _, err = os.Stat(storeDir); os.IsNotExist(err) {
		err = os.MkdirAll(storeDir, 0755)
		if err != nil {
			return
		}
	}

	s.db, err = bolt.Open(s.storePath, 0660)
	if err != nil {
		return
	}

	// Initialize all required buckets
	return s.db.Update(func(tx *bolt.Tx) (err error) {
		s.bucket, err = tx.CreateBucketIfNotExists(s.bucketName)
		return err
	})
}

func (s *boltStore) Close() error {
	return s.db.Close()
}

func (s *boltStore) Put(key string, obj []byte) (err error) {
	m := metrics.GetOrRegisterTimer("fn.store.bolt.Put", nil)
	defer m.UpdateSince(time.Now())

	if !chainstore.IsValidKey(key) {
		return chainstore.ErrInvalidKey
	}
	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)
		return b.Put([]byte(key), obj)
	})
	return
}

func (s *boltStore) Get(key string) (obj []byte, err error) {
	m := metrics.GetOrRegisterTimer("fn.store.bolt.Get", nil)
	defer m.UpdateSince(time.Now())

	if !chainstore.IsValidKey(key) {
		return nil, chainstore.ErrInvalidKey
	}
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)
		obj = b.Get([]byte(key))
		return nil
	})
	return
}

func (s *boltStore) Del(key string) (err error) {
	if !chainstore.IsValidKey(key) {
		return chainstore.ErrInvalidKey
	}
	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)
		return b.Delete([]byte(key))
	})
	return
}
