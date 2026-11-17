// Package store contains pure store implementations and YAML helpers. It must
// not import os, path/filepath or any service package.
package store

import "yukumo-script/core/port"

type (
	Reader           = port.Reader
	Writer           = port.Writer
	ExistenceChecker = port.ExistenceChecker
	Lister           = port.Lister
	Deleter          = port.Deleter
	Store            = port.Store
	Flusher          = port.Flusher
)

// ErrNotFound keeps existing callers on a single sentinel.
var ErrNotFound = port.ErrNotFound
