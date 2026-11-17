package ui

import (
	"fmt"

	"yukumo-script/cli/menu"
	"yukumo-script/cli/prompt"
	"yukumo-script/cli/streams"
)

// Browse runs the interactive menu loop over a menu tree.
func Browse(s streams.Streams, root menu.Node) error {
	stack := []menu.Node{root}
	prompt.Reset(menuParts(stack)...)
	for {
		cur := stack[len(stack)-1]
		children := cur.Children

		labels := make([]string, 0, len(children)+1)
		if len(stack) == 1 {
			labels = append(labels, prompt.LabelExit)
		} else {
			labels = append(labels, prompt.LabelBack)
		}
		for _, c := range children {
			label := c.Name
			if len(c.Children) > 0 {
				label += "/"
			}
			if c.Short != "" {
				label += "  —  " + c.Short
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
		if len(chosen.Children) > 0 {
			stack = append(stack, chosen)
			continue
		}
		prompt.Reset(append(menuParts(stack), chosen.Name)...)
		if chosen.Run == nil {
			_, _ = fmt.Fprintf(s.Err, "command %q is not runnable\n", chosen.Name)
		} else if err := chosen.Run(s); err != nil {
			_, _ = fmt.Fprintln(s.Err, err)
		}
		prompt.Reset(menuParts(stack)...)
	}
}

func menuParts(stack []menu.Node) []string {
	parts := make([]string, 0, len(stack))
	for _, n := range stack {
		parts = append(parts, n.Name)
	}
	return parts
}
