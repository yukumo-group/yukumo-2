// Package cobraui adapts cobra command trees to menu.Node.
package cobraui

import (
	"fmt"

	"yukumo-script/cli/menu"
	"yukumo-script/cli/streams"

	"github.com/spf13/cobra"
)

// FromCobra builds a menu tree from a cobra command tree.
func FromCobra(cmd *cobra.Command) menu.Node {
	children := menuCommands(cmd)
	node := menu.Node{
		Name:  cmd.Name(),
		Short: cmd.Short,
	}
	for _, c := range children {
		node.Children = append(node.Children, FromCobra(c))
	}
	if len(children) == 0 {
		leaf := cmd
		node.Run = func(s streams.Streams) error {
			return runLeaf(leaf, s)
		}
	}
	return node
}

func menuCommands(cmd *cobra.Command) []*cobra.Command {
	var out []*cobra.Command
	for _, c := range cmd.Commands() {
		switch {
		case !c.IsAvailableCommand(), c.Hidden:
			continue
		case c.Name() == "help", c.Name() == "completion", c.Name() == "serve":
			continue
		}
		out = append(out, c)
	}
	return out
}

func runLeaf(cmd *cobra.Command, s streams.Streams) error {
	leaf := &cobra.Command{Use: cmd.Name()}
	leaf.SetIn(s.In)
	leaf.SetOut(s.Out)
	leaf.SetErr(s.Err)
	leaf.SetContext(s.Ctx)
	switch {
	case cmd.RunE != nil:
		return cmd.RunE(leaf, nil)
	case cmd.Run != nil:
		cmd.Run(leaf, nil)
		return nil
	default:
		return fmt.Errorf("command %q is not runnable", cmd.Name())
	}
}
