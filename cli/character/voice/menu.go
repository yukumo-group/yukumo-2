package voice

import (
	"errors"
	"fmt"

	"yukumo-script/cli/prompt"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

// RunMenu is the huh entry under character › edit › voices.
func RunMenu(s streams.Streams, app *api.App, charName string) error {
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
					return runList(s, app, charName, nil)
				})
			case 2:
				runErr = prompt.With("add", func() error {
					return runAdd(s, app, opts, charName, nil)
				})
			case 3:
				runErr = prompt.With("edit", func() error {
					return runEdit(s, app, opts, charName, nil)
				})
			case 4:
				runErr = prompt.With("duplicate", func() error {
					return runDuplicate(s, app, opts, charName, nil)
				})
			case 5:
				runErr = prompt.With("delete", func() error {
					return runDelete(s, app, opts, charName, nil)
				})
			}
			if runErr != nil {
				_, _ = fmt.Fprintln(s.Err, runErr)
			}
		}
	})
}
