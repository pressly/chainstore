package chainstore

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type fs struct {
	storePath string
	perm      os.FileMode // Default: 0755
}

func FileStore(storePath string, perm os.FileMode) (store *fs) {
	if perm == 0 {
		perm = 0755
	}

	store = &fs{storePath: storePath, perm: perm}
	// err = store.Open()
	return
}

func (f *fs) Open() (err error) {
	// Create the path if doesnt exist
	if _, err = os.Stat(f.storePath); os.IsNotExist(err) {
		err = os.MkdirAll(f.storePath, f.perm)
		if err != nil {
			return
		}
	}

	// Check if its a directory and we have rw access
	fd, err := os.Open(f.storePath)
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

func (f *fs) Close() error {
	return nil // noop
}

func (f *fs) Put(key string, obj []byte) (err error) {
	if !IsValidKey(key) {
		return ErrInvalidKey
	}

	if strings.Index(key, "/") > 0 { // folder key
		err = os.MkdirAll(filepath.Dir(filepath.Join(f.storePath, key)), f.perm)
		if err != nil {
			return
		}
	}

	err = ioutil.WriteFile(filepath.Join(f.storePath, key), obj, f.perm)
	return
}

func (f *fs) Get(key string) (obj []byte, err error) {
	if !IsValidKey(key) {
		return nil, ErrInvalidKey
	}

	fp := filepath.Join(f.storePath, key)

	// If the object isn't found, that isn't an error.. just return an empty
	// object.. an error is when we can't talk to the data store
	if _, err = os.Stat(fp); os.IsNotExist(err) {
		return obj, nil
	}

	obj, err = ioutil.ReadFile(fp)
	return
}

func (f *fs) Del(key string) (err error) {
	if string(key[0]) == "/" {
		return ErrInvalidKey
	}
	fp := filepath.Join(f.storePath, key)
	err = os.Remove(fp)
	return
}
