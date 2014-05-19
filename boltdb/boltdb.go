package boltdb

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	. "github.com/nulayer/chainstore"
	"github.com/rcrowley/go-metrics"
)

var lg = log.New(os.Stderr, "", log.LstdFlags)

type Boltdb struct {
	storePath  string
	bucketName []byte

	db     *bolt.DB
	bucket *bolt.Bucket
}

func NewStore(storePath string, bucketName string) (store *Boltdb, err error) {
	store = &Boltdb{storePath: storePath, bucketName: []byte(bucketName)}
	err = store.Open()
	return
}

func (s *Boltdb) Open() (err error) {
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

	if err = s.db.Check(); err != nil {
		if errors, ok := err.(bolt.ErrorList); ok {
			for _, e := range errors {
				lg.Println("[DB ERROR]:", e)
			}

			lg.Println("Deleting bolt db.. and lets carry on")
			os.Remove(s.storePath)
			s.Open()
		}
	}

	// Initialize all required buckets
	return s.db.Update(func(tx *bolt.Tx) (err error) {
		s.bucket, err = tx.CreateBucketIfNotExists(s.bucketName)
		return err
	})
}

func (s *Boltdb) Close() error {
	return s.db.Close()
}

func (s *Boltdb) Put(key string, obj []byte) (err error) {
	m := metrics.GetOrRegisterTimer("fn.store.bolt.Put", nil)
	defer m.UpdateSince(time.Now())

	if !IsValidKey(key) {
		return ErrInvalidKey
	}
	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)
		return b.Put([]byte(key), obj)
	})
	return
}

func (s *Boltdb) Get(key string) (obj []byte, err error) {
	m := metrics.GetOrRegisterTimer("fn.store.bolt.Get", nil)
	defer m.UpdateSince(time.Now())

	if !IsValidKey(key) {
		return nil, ErrInvalidKey
	}
	err = s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)
		obj = b.Get([]byte(key))
		return nil
	})
	return
}

func (s *Boltdb) Del(key string) (err error) {
	if !IsValidKey(key) {
		return ErrInvalidKey
	}
	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.bucketName)
		return b.Delete([]byte(key))
	})
	return
}
