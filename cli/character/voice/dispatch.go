package voice

import (
	"fmt"
	"strings"

	"yukumo-script/cli/streams"
	"yukumo-script/core/api"
)

// Dispatch handles: character <name> voice <verb> ...
func Dispatch(s streams.Streams, app *api.App, charName string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("voice subcommand required (list|get|add|edit|delete|duplicate)")
	}
	verb, rest := args[0], args[1:]
	switch strings.ToLower(verb) {
	case "list":
		return runList(s, app, charName, rest)
	case "get":
		return runGet(s, app, charName, rest)
	case "add":
		return runAdd(s, app, charName, rest)
	case "edit":
		return runEdit(s, app, charName, rest)
	case "delete":
		return runDelete(s, app, charName, rest)
	case "duplicate":
		return runDuplicate(s, app, charName, rest)
	default:
		return fmt.Errorf("unknown voice command %q", verb)
	}
}
