package s3store

import (
	"time"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/nulayer/chainstore"
	"github.com/rcrowley/go-metrics"
)

// TODO: remove go-metrics out of here..

type s3Store struct {
	BucketId, AccessKey, SecretKey string

	conn   *s3.S3
	bucket *s3.Bucket
}

func New(bucketId string, accessKey string, secretKey string) (*s3Store, error) {
	store := &s3Store{BucketId: bucketId, AccessKey: accessKey, SecretKey: secretKey}
	err := store.Open()
	return store, err
}

func (s *s3Store) Open() (err error) {
	auth, err := aws.GetAuth(s.AccessKey, s.SecretKey)
	if err != nil {
		return
	}

	s.conn = s3.New(auth, aws.USEast) // TODO: hardcoded region..?
	s.bucket = s.conn.Bucket(s.BucketId)
	return nil // TODO: no errors ever..?
}

func (s *s3Store) Close() error {
	return nil // TODO: .. nothing to do here..?
}

func (s *s3Store) Put(key string, obj []byte) error {
	m := metrics.GetOrRegisterTimer("fn.store.s3.Put", nil)
	defer m.UpdateSince(time.Now())

	if !chainstore.IsValidKey(key) {
		return chainstore.ErrInvalidKey
	}
	return s.bucket.Put(key, obj, `application/octet-stream`, s3.PublicRead)
}

func (s *s3Store) Get(key string) ([]byte, error) {
	m := metrics.GetOrRegisterTimer("fn.store.s3.Get", nil)
	defer m.UpdateSince(time.Now())

	if !chainstore.IsValidKey(key) {
		return nil, chainstore.ErrInvalidKey
	}

	obj, err := s.bucket.Get(key)
	s3err, _ := err.(*s3.Error)
	if err != nil && s3err.StatusCode != 404 {
		return nil, err
	}
	return obj, nil
}

func (s *s3Store) Del(key string) error {
	if !chainstore.IsValidKey(key) {
		return chainstore.ErrInvalidKey
	}
	return s.bucket.Del(key)
}
