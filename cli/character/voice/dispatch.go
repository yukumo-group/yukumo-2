package voice

import (
	"fmt"
	"strings"

	"yukumo-script/cli/prompt"
	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

type Options struct {
	Interactive bool
}

// Dispatch handles: character <name> voice <verb> ...
func Dispatch(s streams.Streams, app *api.App, opts Options, charName string, args []string) error {
	if len(args) == 0 {
		if !opts.Interactive {
			return fmt.Errorf("voice subcommand required (list|get|add|edit|delete|duplicate)")
		}
		return prompt.With(charName, func() error {
			return RunMenu(s, app, charName)
		})
	}
	verb, rest := args[0], args[1:]
	switch strings.ToLower(verb) {
	case "list":
		return runList(s, app, charName, rest)
	case "get":
		return runGet(s, app, opts, charName, rest)
	case "add":
		return runAdd(s, app, opts, charName, rest)
	case "edit":
		return runEdit(s, app, opts, charName, rest)
	case "delete":
		return runDelete(s, app, opts, charName, rest)
	case "duplicate":
		return runDuplicate(s, app, opts, charName, rest)
	default:
		return fmt.Errorf("unknown voice command %q", verb)
	}
}
