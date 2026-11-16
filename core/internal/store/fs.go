package store

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type FS struct {
	root string
}

func NewFS(root string) *FS {
	return &FS{root: root}
}

func (s *FS) path(name string) string {
	return filepath.Join(s.root, filepath.FromSlash(name))
}

func (s *FS) Read(name string) ([]byte, error) {
	data, err := os.ReadFile(s.path(name))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, ErrNotFound
	}
	return data, err
}

func (s *FS) Write(name string, data []byte) error {
	path := s.path(name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (s *FS) Exists(name string) (bool, error) {
	_, err := os.Stat(s.path(name))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (s *FS) Root() string { return s.root }

func (s *FS) List(dir string) ([]string, error) {
	entries, err := os.ReadDir(s.path(dir))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		out = append(out, e.Name())
	}
	return out, nil
}

func (s *FS) Delete(name string) error {
	err := os.Remove(s.path(name))
	if errors.Is(err, fs.ErrNotExist) {
		return ErrNotFound
	}
	return err
}
