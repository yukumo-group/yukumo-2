package menu_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"yukumo-script/cli/internal/cobraui"
	"yukumo-script/cli/streams"

	"github.com/spf13/cobra"
)

func TestFromCobraFiltersAndRunsLeaf(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	root.AddCommand(&cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Fprint(cmd.OutOrStdout(), "ok")
			return err
		},
	})
	root.AddCommand(&cobra.Command{Use: "hidden", Hidden: true})
	root.AddCommand(&cobra.Command{Use: "help"})
	root.AddCommand(&cobra.Command{Use: "completion"})
	root.AddCommand(&cobra.Command{Use: "serve", Hidden: true})

	spec := cobraui.FromCobra(root)
	if len(spec.Children) != 1 || spec.Children[0].Name != "list" {
		t.Fatalf("children = %+v", spec.Children)
	}
	var buf bytes.Buffer
	s := streams.Streams{Out: &buf, Err: &buf, Ctx: context.Background()}
	if err := spec.Children[0].Run(s); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "ok" {
		t.Fatalf("output = %q, want ok", buf.String())
	}
}
