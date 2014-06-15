package s3store

import (
	"testing"

	"github.com/nulayer/chainstore"
	. "github.com/smartystreets/goconvey/convey"
)

func TestS3Store(t *testing.T) {
	var store chainstore.Store
	var err error

	_ = store
	_ = err

	Convey("S3 Open", t, func() {
		// TODO
	})
}
