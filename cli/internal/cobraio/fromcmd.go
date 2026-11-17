// Package cobraio adapts cobra commands to the CLI Streams contract.
package cobraio

import (
	"context"
	"io"
	"strings"

	"yukumo-script/cli/streams"

	"github.com/spf13/cobra"
)

// FromCmd builds Streams from a cobra command. A nil command yields a
// discard-only Streams so tests and bootstrap code can use the contract
// without a command tree.
func FromCmd(cmd *cobra.Command) streams.Streams {
	if cmd == nil {
		return streams.Streams{
			In:  strings.NewReader(""),
			Out: io.Discard,
			Err: io.Discard,
			Ctx: context.Background(),
		}
	}
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	return streams.Streams{
		In:  cmd.InOrStdin(),
		Out: cmd.OutOrStdout(),
		Err: cmd.ErrOrStderr(),
		Ctx: ctx,
	}
}
