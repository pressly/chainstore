package chainstore

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFsdbStore(t *testing.T) {
	var store Store
	var err error

	Convey("Fsdb Open", t, func() {
		store, err = NewFsdbStore(TempDir(), 0755)
		So(err, ShouldEqual, nil)

		Convey("Put/Get/Del basic data", func() {
			err = store.Put("test.txt", []byte{1, 2, 3, 4})
			So(err, ShouldEqual, nil)

			data, err := store.Get("test.txt")
			So(err, ShouldEqual, nil)
			So(data, ShouldResemble, []byte{1, 2, 3, 4})
		})

		Convey("Disallow invalid keys", func() {
			err = store.Put("test!!!", []byte{1})
			So(err, ShouldEqual, ErrInvalidKey)
		})

		Convey("Auto-creating directories on put", func() {
			err = store.Put("hello/there/everyone.txt", []byte{1, 2, 3, 4})
			So(err, ShouldEqual, nil)
		})

	})
}
