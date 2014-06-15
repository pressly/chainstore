package filestore

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nulayer/chainstore"
)

type fileStore struct {
	storePath string
	perm      os.FileMode // Default: 0755
}

func New(storePath string, perm os.FileMode) *fileStore {
	if perm == 0 {
		perm = 0755
	}

	store := &fileStore{storePath: storePath, perm: perm}
	// err = store.Open()
	return store
}

func (s *fileStore) Open() (err error) {
	// Create the path if doesnt exist
	if _, err = os.Stat(s.storePath); os.IsNotExist(err) {
		err = os.MkdirAll(s.storePath, s.perm)
		if err != nil {
			return
		}
	}

	// Check if its a directory and we have rw access
	fd, err := os.Open(s.storePath)
	if err != nil {
		return
	}
	defer fd.Close()
	fi, err := fd.Stat()
	if err != nil {
		return
	}
	mode := fi.Mode()
	if !mode.IsDir() { // && mode.Perm() // and is rw?
		return errors.New("Store Path is not a directory")
	}
	return
}

func (s *fileStore) Close() error {
	return nil // noop
}

func (s *fileStore) Put(key string, obj []byte) (err error) {
	if !chainstore.IsValidKey(key) {
		return chainstore.ErrInvalidKey
	}

	if strings.Index(key, "/") > 0 { // folder key
		err = os.MkdirAll(filepath.Dir(filepath.Join(s.storePath, key)), s.perm)
		if err != nil {
			return
		}
	}

	err = ioutil.WriteFile(filepath.Join(s.storePath, key), obj, s.perm)
	return
}

func (s *fileStore) Get(key string) (obj []byte, err error) {
	if !chainstore.IsValidKey(key) {
		return nil, chainstore.ErrInvalidKey
	}

	fp := filepath.Join(s.storePath, key)

	// If the object isn't found, that isn't an error.. just return an empty
	// object.. an error is when we can't talk to the data store
	if _, err = os.Stat(fp); os.IsNotExist(err) {
		return obj, nil
	}

	obj, err = ioutil.ReadFile(fp)
	return
}

func (s *fileStore) Del(key string) (err error) {
	if string(key[0]) == "/" {
		return chainstore.ErrInvalidKey
	}
	fp := filepath.Join(s.storePath, key)
	err = os.Remove(fp)
	return
}
