package s3

import (
	"time"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	. "github.com/nulayer/chainstore"
	"github.com/rcrowley/go-metrics"
)

type S3 struct {
	BucketId, AccessKey, SecretKey string

	conn   *s3.S3
	bucket *s3.Bucket
}

func NewStore(bucketId string, accessKey string, secretKey string) (store *S3, err error) {
	store = &S3{BucketId: bucketId, AccessKey: accessKey, SecretKey: secretKey}
	err = store.Open()
	return
}

func (s *S3) Open() (err error) {
	auth, err := aws.GetAuth(s.AccessKey, s.SecretKey)
	if err != nil {
		return
	}

	s.conn = s3.New(auth, aws.USEast) // TODO: hardcoded region..?
	s.bucket = s.conn.Bucket(s.BucketId)
	return nil // TODO: no errors ever..?
}

func (s *S3) Close() error {
	return nil // TODO: .. nothing to do here..?
}

func (s *S3) Put(key string, obj []byte) error {
	m := metrics.GetOrRegisterTimer("fn.store.s3.Put", nil)
	defer m.UpdateSince(time.Now())

	if !IsValidKey(key) {
		return ErrInvalidKey
	}
	// TODO: metadata.........
	return s.bucket.Put(key, obj, `application/octet-stream`, s3.PublicRead)
}

func (s *S3) Get(key string) ([]byte, error) {
	m := metrics.GetOrRegisterTimer("fn.store.s3.Get", nil)
	defer m.UpdateSince(time.Now())

	if !IsValidKey(key) {
		return nil, ErrInvalidKey
	}

	obj, err := s.bucket.Get(key)
	s3err, _ := err.(*s3.Error)
	if err != nil && s3err.StatusCode != 404 {
		return nil, err
	}
	return obj, nil
}

func (s *S3) Del(key string) error {
	if !IsValidKey(key) {
		return ErrInvalidKey
	}
	return s.bucket.Del(key)
}
