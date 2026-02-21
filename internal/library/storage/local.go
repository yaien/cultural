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

// Mount returns the root path of the local storage. The id parameter is ignored since local storage does not require mounting.
func (l *Local) Mount(id string) (string, error) {
	return l.root, nil
}

// Unmount is a no-op for local storage since it does not require unmounting. The id parameter is ignored.
func (l *Local) Unmount(root, id string) error {
	return nil
}
