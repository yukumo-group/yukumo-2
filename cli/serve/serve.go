package serve

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"yukumo-script/core/api/rpc"
)

func Loop(ctx context.Context, reg *rpc.Registry, in io.Reader, out io.Writer, errOut io.Writer) error {
	if reg == nil {
		return fmt.Errorf("serve: registry is nil")
	}
	r := bufio.NewReader(in)
	w := bufio.NewWriter(out)
	for {
		line, err := r.ReadBytes('\n')
		if len(line) > 0 {
			payload := trimNL(line)
			if len(payload) > 0 {
				resp := reg.Handle(ctx, payload)
				if len(resp) > 0 {
					if _, werr := w.Write(resp); werr != nil {
						return werr
					}
					if werr := w.WriteByte('\n'); werr != nil {
						return werr
					}
					if werr := w.Flush(); werr != nil {
						return werr
					}
				}
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			if errOut != nil {
				fmt.Fprintln(errOut, err)
			}
			return err
		}
	}
}

func trimNL(b []byte) []byte {
	for len(b) > 0 && (b[len(b)-1] == '\n' || b[len(b)-1] == '\r') {
		b = b[:len(b)-1]
	}
	return b
}
