// Package streams defines the IO contract shared by the CLI shell, ui, and
// render layers. It intentionally has no dependency on cobra.
package streams

import (
	"context"
	"io"
)

// Streams carries the standard input, output, error, and context for a
// command invocation.
type Streams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
	Ctx context.Context
}
