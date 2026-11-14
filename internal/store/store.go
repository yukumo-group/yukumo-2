package store

import "errors"

var ErrNotFound = errors.New("store: not found")

type Reader interface {
	Read(name string) ([]byte, error)
}

type Writer interface {
	Write(name string, data []byte) error
}

type ExistenceChecker interface {
	Exists(name string) (bool, error)
}

type Lister interface {
	List(dir string) ([]string, error)
}

type Deleter interface {
	Delete(name string) error
}

type Store interface {
	Reader
	Writer
	ExistenceChecker
	Lister
	Deleter
}
