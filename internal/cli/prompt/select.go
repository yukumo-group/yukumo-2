package prompt

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

const (
	LabelBack   = "←  Back"
	LabelExit   = "←  Exit"
	LabelParent = "←  Parent directory"
	LabelSave   = "✓  Save"
	LabelCancel = "←  Cancel"
)

func Select(segment, desc string, labels []string, defaultIdx int) (int, error) {
	if len(labels) == 0 {
		return 0, fmt.Errorf("no options")
	}
	options := make([]huh.Option[int], len(labels))
	for i, l := range labels {
		options[i] = huh.NewOption(l, i)
	}
	idx := defaultIdx
	if idx < 0 || idx >= len(labels) {
		idx = 0
	}
	sel := huh.NewSelect[int]().
		Title(Title(segment)).
		Options(options...).
		Value(&idx)
	if desc != "" {
		sel = sel.Description(desc)
	}
	if err := run(sel); err != nil {
		return 0, err
	}
	return idx, nil
}
