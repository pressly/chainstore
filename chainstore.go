package chainstore

import (
	"errors"
	"io/ioutil"
	"regexp"
)

var (
	ErrInvalidKey = errors.New("Invalid key")
	ErrNoStores   = errors.New("No stores have been provided to chain")
)

const (
	MAX_KEY_LEN = 256
)

type Store interface {
	Open() error
	Close() error
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
	Del(key string) error
}

func New(stores ...Store) (Store, error) {
	return NewChain(stores...)
}

//--

func IsValidKey(key string) bool {
	// TODO: should this regexp be prebuilt..?
	m, _ := regexp.MatchString(`(i?)[^a-z0-9\/_\-:\.]`, key)
	return !m && len(key) <= MAX_KEY_LEN
}

func TempDir() string {
	path, _ := ioutil.TempDir("", "chainstore-")
	return path
}
