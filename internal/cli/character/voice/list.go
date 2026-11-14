package voice

import (
	svccharacter "yukumo-script/internal/service/character"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func runList(cmd *cobra.Command, svc *svccharacter.Service, charName string, args []string) error {
	if len(args) != 0 {
		return cmd.Usage()
	}
	voices, err := svc.ListVoices(charName)
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent(2)
	defer enc.Close()
	return enc.Encode(voices)
}
