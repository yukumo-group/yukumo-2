package cli

import (
	"fmt"
	"strings"

	"yukumo-script/internal/cli/prompt"

	"github.com/spf13/cobra"
)

func browseCommands(start *cobra.Command) error {
	stack := []*cobra.Command{start}
	for {
		cur := stack[len(stack)-1]
		children := menuCommands(cur)

		labels := make([]string, 0, len(children)+1)
		if len(stack) == 1 {
			labels = append(labels, prompt.LabelExit)
		} else {
			labels = append(labels, prompt.LabelBack)
		}
		for _, c := range children {
			label := c.Name()
			if len(menuCommands(c)) > 0 {
				label += "/"
			}
			if s := c.Short; s != "" {
				label += "  —  " + s
			}
			labels = append(labels, label)
		}

		idx, err := prompt.Select(menuTitle(stack), "↑↓ move · enter", labels, 0)
		if err != nil {
			return err
		}
		if idx == 0 {
			if len(stack) == 1 {
				return nil
			}
			stack = stack[:len(stack)-1]
			continue
		}

		chosen := children[idx-1]
		if len(menuCommands(chosen)) > 0 {
			stack = append(stack, chosen)
			continue
		}
		if err := runCommand(chosen); err != nil {
			fmt.Fprintln(chosen.ErrOrStderr(), err)
			continue
		}
	}
}

func menuCommands(cmd *cobra.Command) []*cobra.Command {
	var out []*cobra.Command
	for _, c := range cmd.Commands() {
		switch {
		case !c.IsAvailableCommand(), c.Hidden:
			continue
		case c.Name() == "help", c.Name() == "completion":
			continue
		}
		out = append(out, c)
	}
	return out
}

func menuTitle(stack []*cobra.Command) string {
	parts := make([]string, 0, len(stack))
	for _, c := range stack {
		parts = append(parts, c.Name())
	}
	return strings.Join(parts, " › ")
}

func runCommand(cmd *cobra.Command) error {
	switch {
	case cmd.RunE != nil:
		return cmd.RunE(cmd, nil)
	case cmd.Run != nil:
		cmd.Run(cmd, nil)
		return nil
	default:
		return fmt.Errorf("command %q is not runnable", cmd.Name())
	}
}

func enableInteractiveBrowse(root *cobra.Command) {
	attachBrowse(root)
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		for _, c := range menuCommands(cmd) {
			if len(menuCommands(c)) > 0 && c.Run == nil && c.RunE == nil {
				attachBrowse(c)
			}
			walk(c)
		}
	}
	walk(root)
}

func attachBrowse(cmd *cobra.Command) {
	start := cmd
	cmd.RunE = func(*cobra.Command, []string) error {
		return browseCommands(start)
	}
}
