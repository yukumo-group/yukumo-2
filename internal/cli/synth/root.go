package synth

import (
	"fmt"
	"strings"

	"yukumo-script/internal/cli/prompt"
	svccharacter "yukumo-script/internal/service/character"
	svcvoice "yukumo-script/internal/service/character/voice"
	svcsynth "yukumo-script/internal/service/synth"

	"github.com/spf13/cobra"
)

type Options struct {
	Interactive bool
}

func NewCmd(
	synthSvc *svcsynth.Service,
	chars *svccharacter.Service,
	voices *svcvoice.Service,
	opts Options,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "synth",
		Short: "Synthesize speech from a character voice",
		Long: `Synthesize speech.

  synth -v character.voice [-V 1|2|10] [-s factor] [-o path] [-p] TEXT
  synth custom   (stub; not implemented yet)`,
		Args:                  cobra.ArbitraryArgs,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSynth(cmd, synthSvc, chars, voices, opts, args)
		},
	}

	cmd.Flags().StringP("voice", "v", "", "Voice ref character.voice")
	cmd.Flags().IntP("version", "V", 0, "Engine version 1|2|10 (default: first available)")
	cmd.Flags().Float64P("speed", "s", 1, "Speed factor multiplied with voice speed")
	cmd.Flags().StringP("output", "o", "", "Output WAV path")
	cmd.Flags().BoolP("preview", "p", false, "Play synthesized audio")

	cmd.AddCommand(newVoiceCmd(synthSvc, chars, voices, opts))
	cmd.AddCommand(newCustomCmd())
	return cmd
}

func newCustomCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "custom",
		Short: "Synthesize with fully custom params (not implemented)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("synth custom is not implemented yet")
		},
	}
}

func runSynth(
	cmd *cobra.Command,
	synthSvc *svcsynth.Service,
	chars *svccharacter.Service,
	voices *svcvoice.Service,
	opts Options,
	args []string,
) error {
	voiceRef, _ := cmd.Flags().GetString("voice")
	version, _ := cmd.Flags().GetInt("version")
	speed, _ := cmd.Flags().GetFloat64("speed")
	output, _ := cmd.Flags().GetString("output")
	preview, _ := cmd.Flags().GetBool("preview")

	text := strings.TrimSpace(strings.Join(args, " "))
	if err := prompt.EnsureString("text", &text); err != nil {
		return err
	}
	if text == "" {
		return fmt.Errorf("text is required")
	}
	if err := ensureVoiceRef(chars, voices, opts.Interactive, &voiceRef); err != nil {
		return err
	}

	_, err := runRequest(cmd, synthSvc, svcsynth.Request{
		VoiceRef:    voiceRef,
		Text:        text,
		Version:     version,
		SpeedFactor: speed,
		Output:      output,
		Preview:     preview,
	}, false)
	return err
}
