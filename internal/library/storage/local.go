package storage

import (
	"io"
	"os"
)

var _ Storage = (*Local)(nil)

type Local struct {
	root string
}

func NewLocal(root string) *Local {
	return &Local{
		root: root,
	}
}

func (l *Local) Put(id string, size int64, data io.Reader) error {
	root, err := os.OpenRoot(l.root)
	if err != nil {
		return err
	}
	defer root.Close()

	file, err := root.Create(id)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.CopyN(file, data, size)
	if err != nil {
		return err
	}

	return nil
}

func (l *Local) Remove(id string) error {
	root, err := os.OpenRoot(l.root)
	if err != nil {
		return err
	}
	defer root.Close()

	return root.Remove(id)
}

func (l *Local) Get(id string) (io.ReadCloser, error) {
	return os.OpenInRoot(l.root, id)
}
