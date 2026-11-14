package voice

import (
	"errors"
	"fmt"

	"yukumo-script/internal/cli/prompt"
	svcvoice "yukumo-script/internal/service/character/voice"

	"github.com/spf13/cobra"
)

// RunMenu is the huh entry under character › edit › voices.
func RunMenu(cmd *cobra.Command, svc *svcvoice.Service, charName string) error {
	opts := Options{Interactive: true}
	return prompt.With("voices", func() error {
		for {
			labels := []string{
				prompt.LabelBack,
				"list  —  Print voices",
				"add  —  Create a voice",
				"edit  —  Update a voice",
				"duplicate  —  Copy a voice",
				"delete  —  Remove a voice",
			}
			idx, err := prompt.Select("", "↑↓ move · enter", labels, 0)
			if err != nil {
				if errors.Is(err, prompt.ErrAborted) {
					return nil
				}
				return err
			}
			var runErr error
			switch idx {
			case 0:
				return nil
			case 1:
				runErr = prompt.With("list", func() error {
					return runList(cmd, svc, charName, nil)
				})
			case 2:
				runErr = prompt.With("add", func() error {
					return runAdd(cmd, svc, opts, charName, nil)
				})
			case 3:
				runErr = prompt.With("edit", func() error {
					return runEdit(cmd, svc, opts, charName, nil)
				})
			case 4:
				runErr = prompt.With("duplicate", func() error {
					return runDuplicate(cmd, svc, opts, charName, nil)
				})
			case 5:
				runErr = prompt.With("delete", func() error {
					return runDelete(cmd, svc, opts, charName, nil)
				})
			}
			if runErr != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), runErr)
			}
		}
	})
}
