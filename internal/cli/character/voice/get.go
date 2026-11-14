package voice

import (
	svcvoice "yukumo-script/internal/service/character/voice"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func runGet(cmd *cobra.Command, svc *svcvoice.Service, opts Options, charName string, args []string) error {
	name := ""
	if len(args) >= 1 {
		name = args[0]
	}
	if opts.Interactive || name == "" {
		if err := ensureVoiceName(svc, charName, &name); err != nil {
			return err
		}
	}
	if name == "" {
		return cmd.Usage()
	}
	v, err := svc.Get(charName, name)
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent(2)
	defer enc.Close()
	return enc.Encode(v)
}
