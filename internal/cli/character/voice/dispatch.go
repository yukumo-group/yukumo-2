package voice

import (
	"fmt"
	"strings"

	"yukumo-script/internal/cli/prompt"
	svcvoice "yukumo-script/internal/service/character/voice"

	"github.com/spf13/cobra"
)

type Options struct {
	Interactive bool
}

// Dispatch handles: character <name> voice <verb> ...
func Dispatch(cmd *cobra.Command, svc *svcvoice.Service, opts Options, charName string, args []string) error {
	if len(args) == 0 {
		if !opts.Interactive {
			return fmt.Errorf("voice subcommand required (list|get|add|edit|delete|duplicate)")
		}
		return prompt.With(charName, func() error {
			return RunMenu(cmd, svc, charName)
		})
	}
	verb, rest := args[0], args[1:]
	switch strings.ToLower(verb) {
	case "list":
		return runList(cmd, svc, charName, rest)
	case "get":
		return runGet(cmd, svc, opts, charName, rest)
	case "add":
		return runAdd(cmd, svc, opts, charName, rest)
	case "edit":
		return runEdit(cmd, svc, opts, charName, rest)
	case "delete":
		return runDelete(cmd, svc, opts, charName, rest)
	case "duplicate":
		return runDuplicate(cmd, svc, opts, charName, rest)
	default:
		return fmt.Errorf("unknown voice command %q", verb)
	}
}
