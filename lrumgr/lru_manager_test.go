package lrumgr

import (
	"testing"

	"github.com/nulayer/chainstore"
	"github.com/nulayer/chainstore/filestore"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLRUManager(t *testing.T) {
	var err error
	var store chainstore.Store
	var lru *LruManager
	var capacity int64 = 20

	Convey("LRUManager", t, func() {
		storeDir := chainstore.TempDir()

		store = filestore.New(storeDir, 0755)
		lru = New(capacity, store)

		// based on 10% cushion
		lru.Put("peter", []byte{1, 2, 3})
		lru.Put("jeff", []byte{4})
		lru.Put("julia", []byte{5, 6, 7, 8, 9, 10})
		lru.Put("janet", []byte{11, 12, 13})
		lru.Put("ted", []byte{14, 15, 16, 17, 18})

		remaining := capacity - 18
		So(lru.Capacity(), ShouldEqual, remaining)

		remaining = remaining + 4
		err = lru.Put("agnes", []byte{20, 21, 22, 23, 24, 25})
		So(lru.Capacity(), ShouldEqual, remaining)
		So(err, ShouldEqual, nil)
	})
}
