package main

import (
	"testing"
	"time"

	"github.com/nulayer/chainstore"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBasicChain(t *testing.T) {
	var s1, s2, chain chainstore.Store
	var err error

	Convey("Chain", t, func() {
		storeDir := chainstore.TempDir()
		err = nil

		s1 = chainstore.MemCacheStore(10) // 10 byte capacity
		So(err, ShouldEqual, nil)

		s2 = chainstore.FileStore(storeDir+"/s2", 0755)
		So(err, ShouldEqual, nil)

		chain, err = chainstore.New(s1, s2)
		So(err, ShouldEqual, nil)

		Convey("Put on chain", func() {
			err = chain.Put("hi", []byte{1, 2, 3, 4})
			So(err, ShouldEqual, nil)
			time.Sleep(10e6) // wait..

			v, err := s1.Get("hi")
			So(err, ShouldEqual, nil)
			So(v, ShouldResemble, []byte{1, 2, 3, 4})

			v, err = s2.Get("hi")
			So(err, ShouldEqual, nil)
			So(v, ShouldResemble, []byte{1, 2, 3, 4})
		})
	})
}

func TestMiddlewareChain(t *testing.T) {
	Convey("Chain with middleware", t, func() {
		chain, err := chainstore.New(
			chainstore.MemCacheStore(100),
			chainstore.FileStore(chainstore.TempDir(), 0755),
		)

		So(err, ShouldEqual, nil)

		v, err := chain.Get("x")
		So(len(v), ShouldEqual, 0)

		err = chain.Put("x", []byte("hiiiiiii!"))
		So(err, ShouldEqual, nil)

		time.Sleep(10e6) // wait..

		v, err = chain.Get("x")
		So(len(v), ShouldBeGreaterThan, 0)
	})

	// HRMM... where can I go to find other examples...?

	Convey("More complicated middleware defs", t, func() {
		chain, err := chainstore.New(
			chainstore.MemCacheStore(100),
			chainstore.FileStore(chainstore.TempDir(), 0755),
		)

		So(err, ShouldEqual, nil)

		v, err := chain.Get("x")
		So(len(v), ShouldEqual, 0)

		err = chain.Put("x", []byte("hiiiiiii!"))
		So(err, ShouldEqual, nil)

		time.Sleep(10e6) // wait..

		v, err = chain.Get("x")
		So(len(v), ShouldBeGreaterThan, 0)
	})
}

/*
func TestMiddlewareChain(t *testing.T) {
  // var s1, s2, chain Store
  // var err error

  Convey("Chain with middleware", t, func() {
    storeDir := TempDir()

    s1 = MemCacheStore(100)
    lru := LRUManager(100)
    chain, err = chainstore.New(lru, s1)
    So(err, ShouldEqual, nil)
  })

  Convey("Chain with middleware, take 2", t, func() {
    storeDir := TempDir()

    // This will create a &memCache{}. It does break the naming
    // convention of NewX(), as well, it doesn't return any errors
    // here so they can be defined inline the chainstore.New() call
    s1 = MemCacheStore(100)

    // BUT... what do we do with the errors then?
    // one way, we move all error handling of a store to its Open()
    // method, which returns errors, so a store cannot be used until
    // its Opened(). chainstore.New(...), will Open() each store,
    // but what if the store had already been Open()'d?

    // Chainstore definition example:
    metricsMgr = metrics.NewManger(&m)

    chain, err = chainstore.New(
      logger.NewManager(l),
      LRUManager(100),
      metricsMgr,
      s1,
      metricsMgr,
      boltdb.NewStore("/tmp/store.db"),
      s3.NewStore("key", "x"))

    // Another idea: a chainstore provides a channel for errors
    // that occur through the chain. It might not be ideal, but
    // there are parts of the chain that happen asynchronously
    chain.Errors(func(err error) {
    })

    // Finally, after seeing this, perhaps we should have nested chains:
    c, _ = chainstore.New(
      logger.NewManager(l),
      chainstore.New(metricsMgr, LRUManager(100), FileStore("/tmp", 0755)),
      chainstore.New(AsyncManager, boltdb.NewStore("/tmp/s.db").Use(metricsMgr), s3.NewStore("x", "y"))
    )

  })
}
*/
