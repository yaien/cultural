package storage

import "io"

type Storage interface {
	Put(id string, size int64, data io.Reader) error
	Remove(id string) error
	Get(id string) (io.ReadCloser, error)
}
