// Package port contains the public contracts between the application core and
// its hosts (CLI filesystem, browser wasm, tests).
//
// Keys are logical: slash-separated, relative, never absolute and never
// escaping the adapter root. Implementations own platform details such as
// directories, IndexedDB object stores or in-memory maps.
package port

import (
	"context"
	"errors"
)

// ErrNotFound is returned by store reads and deletes for a missing key.
var ErrNotFound = errors.New("store: not found")

type Reader interface {
	Read(key string) ([]byte, error)
}

type Writer interface {
	Write(key string, data []byte) error
}

type ExistenceChecker interface {
	Exists(key string) (bool, error)
}

// Lister returns immediate child keys under prefix, sorted.
type Lister interface {
	List(prefix string) ([]string, error)
}

type Deleter interface {
	Delete(key string) error
}

// Store is a byte-level key/value store. Reads and writes are synchronous; an
// adapter may persist asynchronously behind its own cache and Flush.
type Store interface {
	Reader
	Writer
	ExistenceChecker
	Lister
	Deleter
}

// Flusher is optionally implemented by stores that keep pending writes.
type Flusher interface {
	Flush(ctx context.Context) error
}
