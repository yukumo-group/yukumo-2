// Package fs provides the filesystem-backed port implementations. It is the
// only store package allowed to use os and path/filepath.
package fs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"yukumo-script/core/port"
)

// Store is a port.Store rooted at a directory. It also implements
// port.ArtifactStore.
type Store struct {
	root string
}

func New(root string) *Store {
	return &Store{root: root}
}

func (s *Store) Root() string { return s.root }

func cleanKey(key string) (string, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return "", fmt.Errorf("store: empty key")
	}
	if strings.HasPrefix(key, "/") || filepath.IsAbs(key) || strings.Contains(key, `\`) {
		return "", fmt.Errorf("store: invalid key %q", key)
	}
	clean := filepath.ToSlash(filepath.Clean(key))
	if clean == "." || strings.HasPrefix(clean, "../") || clean == ".." {
		return "", fmt.Errorf("store: key escapes root: %q", key)
	}
	return clean, nil
}

func (s *Store) path(key string) (string, error) {
	clean, err := cleanKey(key)
	if err != nil {
		return "", err
	}
	return filepath.Join(s.root, filepath.FromSlash(clean)), nil
}

func (s *Store) Read(key string) ([]byte, error) {
	path, err := s.path(key)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, port.ErrNotFound
	}
	return data, err
}

func (s *Store) Write(key string, data []byte) error {
	path, err := s.path(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (s *Store) Exists(key string) (bool, error) {
	path, err := s.path(key)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (s *Store) List(prefix string) ([]string, error) {
	dir, err := s.path(prefix)
	if err != nil {
		if prefix == "" {
			dir = s.root
		} else {
			return nil, err
		}
	}
	entries, err := os.ReadDir(dir)
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
	sort.Strings(out)
	return out, nil
}

func (s *Store) Delete(key string) error {
	path, err := s.path(key)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if errors.Is(err, fs.ErrNotExist) {
		return port.ErrNotFound
	}
	return err
}

// Path is display-only.
func (s *Store) Path(key string) string {
	path, err := s.path(key)
	if err != nil {
		return ""
	}
	return path
}

// Output writes explicit -o paths.
type Output struct{}

func (Output) Write(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
