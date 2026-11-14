package voice

import (
	"fmt"
	"strings"

	svccharacter "yukumo-script/internal/service/character"

	"github.com/spf13/cobra"
)

type Options struct {
	Interactive bool
}

// Dispatch handles: character <name> voice <verb> ...
func Dispatch(cmd *cobra.Command, svc *svccharacter.Service, opts Options, charName string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("voice subcommand required (list|get|add|edit|delete)")
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
	default:
		return fmt.Errorf("unknown voice command %q", verb)
	}
}
