package chainstoretest

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/nulayer/chainstore"
	"github.com/nulayer/chainstore/filestore"
	"github.com/nulayer/chainstore/logmgr"
	"github.com/nulayer/chainstore/memstore"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBasicChain(t *testing.T) {
	var ms, fs, chain chainstore.Store
	var err error

	logger := log.New(os.Stdout, "", log.LstdFlags)

	Convey("Chain", t, func() {
		storeDir := chainstore.TempDir()
		err = nil

		ms = memstore.New(100)
		fs = filestore.New(storeDir+"/filestore", 0755)

		chain = chainstore.New(
			logmgr.New(logger, ""),
			ms,
			chainstore.Async(
				logmgr.New(logger, "async"),
				fs,
			),
		)
		err = chain.Open()
		So(err, ShouldEqual, nil)

		Convey("Put", func() {
			v := []byte("value")
			err = chain.Put("k", v)
			So(err, ShouldEqual, nil)

			val, err := chain.Get("k")
			So(err, ShouldEqual, nil)
			So(v, ShouldResemble, v)

			val, err = ms.Get("k")
			So(err, ShouldEqual, nil)
			So(val, ShouldResemble, v)

			time.Sleep(10e6) // wait for async operation..

			val, err = fs.Get("k")
			So(err, ShouldEqual, nil)
			So(val, ShouldResemble, v)
		})

	})

}

/*

c := chainstore.New(
	logger,
	memstore,
	filestore,
	bg_manager,
	boltdb,
	s3
)

c := chainstore.New(
	logger,
	chainstore.New(memstore, s2),
	chainstore.New(filestore, s3store).async()
)



... okay.. something like that is okay for normal usage, with go support..

* how do we get to support metrics tho..? which wraps the time
same example:


this will work:..

c := chainstore.New(
	logger,
	chainstore.New(memstore, s2),
	chainstore.Async(metricsmgr.New(chainstore.New(filestore, s3store), "blah"))
)

... a bit crazy.. now lets add......

* batch > lru > bolt

c := chainstore.New(
	logmgr.New(l, ""),
	memstore.New(1000),
	chainstore.Async(
		logmgr.New(l, "async"),
		metricsmgr.New(
			"bolt", &metricsmgr.Config{a, b, c},
			batchmgr.New(10),
			lrumgr.New(5000, boltstore.New("/tmp/bolt.db", 0755)),
		),
		metricsmgr.New(
			"s3", &metricsmgr.Config{a, b, c}
			s3store.New("bucket", "u", "p")
		)
	)
)


























































*/
