package voice

import (
	"fmt"
	"strings"

	"yukumo-script/cli/prompt"
	"yukumo-script/cli/streams"
	"yukumo-script/cli/ui"
	"yukumo-script/core/api"
)

// Dispatch handles: character <name> voice <verb> ...
func Dispatch(s streams.Streams, app *api.App, interactive bool, charName string, args []string) error {
	if len(args) == 0 {
		if !interactive {
			return fmt.Errorf("voice subcommand required (list|get|add|edit|delete|duplicate)")
		}
		return prompt.With(charName, func() error {
			return ui.VoiceMenu(s, app, charName)
		})
	}
	verb, rest := args[0], args[1:]
	switch strings.ToLower(verb) {
	case "list":
		return runList(s, app, charName, rest)
	case "get":
		return runGet(s, app, interactive, charName, rest)
	case "add":
		return runAdd(s, app, interactive, charName, rest)
	case "edit":
		return runEdit(s, app, interactive, charName, rest)
	case "delete":
		return runDelete(s, app, interactive, charName, rest)
	case "duplicate":
		return runDuplicate(s, app, interactive, charName, rest)
	default:
		return fmt.Errorf("unknown voice command %q", verb)
	}
}
