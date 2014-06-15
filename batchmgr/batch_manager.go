package batchmgr

import "github.com/nulayer/chainstore"

// import "container/list"

type batchManager struct {
	num   int64
	batch map[string][]byte // todo, should be array.. pop on/off..
	// list  *list.List
}

func New(num int64) *batchManager {
	return &batchManager{num, make(map[string][]byte, num)}
}

func (b *batchManager) Open() (err error) { return }

func (b *batchManager) Close() (err error) { return }

func (b *batchManager) Put(key string, obj []byte) (err error) {
	if !chainstore.IsValidKey(key) {
		return chainstore.ErrInvalidKey
	}
	return
}

func (b *batchManager) Get(key string) (obj []byte, err error) {
	return
}

func (b *batchManager) Del(key string) (err error) {
	return
}
