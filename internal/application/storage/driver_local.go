package storage

import (
	"fmt"
	"io"
	"os"
	"path"
)

var _ Driver = (*Local)(nil)

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

	if _, err = io.CopyN(file, data, size); err != nil {
		return fmt.Errorf("failed in copy data to file: %w", err)
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

func (l *Local) Mount(id string) (dir, src string, err error) {
	return l.root, path.Join(l.root, id), nil
}

func (l *Local) Unmount(dir string) error {
	return nil
}
