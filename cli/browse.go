package cli

import (
	"fmt"

	"yukumo-script/cli/prompt"

	"github.com/spf13/cobra"
)

func browseCommands(start *cobra.Command) error {
	stack := []*cobra.Command{start}
	prompt.Reset(menuParts(stack)...)
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

		prompt.Reset(menuParts(stack)...)
		idx, err := prompt.Select("", "↑↓ move · enter", labels, 0)
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
		prompt.Reset(append(menuParts(stack), chosen.Name())...)
		if err := runCommand(chosen); err != nil {
			_, _ = fmt.Fprintln(chosen.ErrOrStderr(), err)
		}
		prompt.Reset(menuParts(stack)...)
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

func menuParts(stack []*cobra.Command) []string {
	parts := make([]string, 0, len(stack))
	for _, c := range stack {
		parts = append(parts, c.Name())
	}
	return parts
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
