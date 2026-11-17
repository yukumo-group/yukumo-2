package store

import (
	"context"
	"path"
	"sort"
	"strings"
	"sync"
)

// Memory is a pure in-memory store. It also implements port.ArtifactStore and
// port.Flusher so tests and wasm hosts can use it directly.
type Memory struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewMemory() *Memory {
	return &Memory{data: make(map[string][]byte)}
}

func (m *Memory) Read(name string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.data[name]
	if !ok {
		return nil, ErrNotFound
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out, nil
}

func (m *Memory) Write(name string, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cp := make([]byte, len(data))
	copy(cp, data)
	m.data[name] = cp
	return nil
}

func (m *Memory) Exists(name string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.data[name]
	return ok, nil
}

func (m *Memory) List(dir string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	prefix := strings.Trim(path.Clean("/"+dir), "/")
	if prefix == "." {
		prefix = ""
	}
	var out []string
	for name := range m.data {
		if prefix == "" {
			if !strings.Contains(name, "/") {
				out = append(out, name)
			}
			continue
		}
		if strings.HasPrefix(name, prefix+"/") {
			rest := name[len(prefix)+1:]
			if !strings.Contains(rest, "/") {
				out = append(out, rest)
			}
		}
	}
	sort.Strings(out)
	return out, nil
}

func (m *Memory) Delete(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.data[name]; !ok {
		return ErrNotFound
	}
	delete(m.data, name)
	return nil
}

// Path makes Memory usable as a port.ArtifactStore; it is display-only.
func (m *Memory) Path(key string) string { return key }

// Flush makes Memory usable as a port.Flusher.
func (m *Memory) Flush(context.Context) error { return nil }
