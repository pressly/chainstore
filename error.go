package chainstore

import (
	"errors"
)

var (
	ErrInvalidKey    = errors.New("Invalid key.")
	ErrMissingStores = errors.New("No stores provided.")
	ErrNoSuchKey     = errors.New("No such key.")
)
