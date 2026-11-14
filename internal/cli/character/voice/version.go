package voice

import (
	"fmt"

	"yukumo-script/internal/cli/prompt"
	"yukumo-script/internal/domain"
)

func SelectVersion(out *int) error {
	if out == nil {
		return fmt.Errorf("version is nil")
	}
	labels := []string{"1 (AquesTalk1)", "2 (AquesTalk2)", "10 (AquesTalk10)"}
	values := []int{domain.Version1, domain.Version2, domain.Version10}
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
