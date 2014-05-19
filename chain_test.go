package chainstore

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChain(t *testing.T) {
	var s1, s2, chain Store
	var err error

	Convey("Chain", t, func() {
		storeDir := TempDir()

		s1, err = NewMemCacheStore(10) // 10 byte capacity
		So(err, ShouldEqual, nil)

		s2, err = NewFsdbStore(storeDir+"/s2", 0755)
		So(err, ShouldEqual, nil)

		chain, err = NewChain(s1, s2)
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
