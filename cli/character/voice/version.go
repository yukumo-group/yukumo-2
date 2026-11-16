package voice

import (
	"fmt"

	"yukumo-script/cli/prompt"
	"yukumo-script/core/api"
)

func SelectVersion(out *int) error {
	if out == nil {
		return fmt.Errorf("version is nil")
	}
	labels := []string{"1 (AquesTalk1)", "2 (AquesTalk2)", "10 (AquesTalk10)"}
	values := []int{api.Version1, api.Version2, api.Version10}
	focus := 0
	for i, v := range values {
		if v == *out {
			focus = i
			break
		}
	}
	idx, err := prompt.Select("version", "↑↓ move · enter", labels, focus)
	if err != nil {
		return err
	}
	*out = values[idx]
	return nil
}
