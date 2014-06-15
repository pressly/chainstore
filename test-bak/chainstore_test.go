package chainstoretest

import (
	"testing"
	"time"

	"github.com/nulayer/chainstore"
	"github.com/nulayer/chainstore/filestore"
	"github.com/nulayer/chainstore/memstore"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBasicChain(t *testing.T) {
	var s1, s2, chain chainstore.Store
	var err error

	Convey("Chain", t, func() {
		storeDir := chainstore.TempDir()
		err = nil

		s1 = memstore.New(10) // 10 byte capacity
		So(err, ShouldEqual, nil)

		s2 = filestore.New(storeDir+"/s2", 0755)
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

/*
func TestMiddlewareChain(t *testing.T) {
	// var s1, s2, chain Store
	// var err error

	Convey("Chain with middleware", t, func() {
		storeDir := TempDir()

		s1 = memstore.New(100)
		lru := lrumgr.New(100)
		chain, err = chainstore.New(lru, s1)
		So(err, ShouldEqual, nil)
	})

	Convey("Chain with middleware, take 2", t, func() {
		storeDir := TempDir()

		// This will create a &memCache{}. It does break the naming
		// convention of NewX(), as well, it doesn't return any errors
		// here so they can be defined inline the chainstore.New() call
		s1 = memstore.New(100)

		// BUT... what do we do with the errors then?
		// one way, we move all error handling of a store to its Open()
		// method, which returns errors, so a store cannot be used until
		// its Opened(). chainstore.New(...), will Open() each store,
		// but what if the store had already been Open()'d?

		// Chainstore definition example:
		metricsMgr = metricsmgr.New(&m)

		chainboltdb.NewStore

		levelstore.New()
		boltstore
		metricmgr
		s3store
		batchmgr
		asyncmgr or bgmgr
		logmgr.New
		filestore.New
		batchmgr.New
		memstore.New
		lrumgr.New

		cstore.New() ...?

		.Open, .Close, .Put, .Get, .Del

		import (
			"github.com/nulayer/chainstore"
			"github.com/nulayer/chainstore/batchmgr"
			"github.com/nulayer/chainstore/memstore"
			"github.com/nulayer/chainstore/boltstore"
			"github.com/nulayer/chainstore/metricsmgr"
			"github.com/nulayer/chainstore/s3store"
		)

		chain, err = chainstore.New(
			logmgr.New(l),
			lrumgr.New(100),
			metricsMgr,
			s1,
			metricsMgr,
			boltstore.New("/tmp/store.db"),
			s3store.New("key", "x"))

		// Another idea: a chainstore provides a channel for errors
		// that occur through the chain. It might not be ideal, but
		// there are parts of the chain that happen asynchronously
		chain.Errors(func(err error) {
		})

		// Finally, after seeing this, perhaps we should have nested chains:
		c, _ = chainstore.New(
			logger.NewManager(l),
			chainstore.New(metricsMgr, lrumgr.New(100), filestore.New("/tmp", 0755)),
			chainstore.New(bgmgr.New(), boltstore.New("/tmp/s.db").Use(metricsMgr), s3store.New("x", "y"))
		)

	})
}
*/
