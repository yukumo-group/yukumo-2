package synth

import (
	"fmt"
	"strings"

	"yukumo-script/cli/prompt"
	"yukumo-script/core/api"

	"github.com/spf13/cobra"
)

type Options struct {
	Interactive bool
}

func NewCmd(app *api.App, opts Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "synth",
		Short: "Synthesize speech from a character voice",
		Long: `Synthesize speech.

  synth -v character:voice [-V 1|2|10] [-s factor] [-o path] [-p] [--reuse-cache] TEXT
  synth -P /path/character.yaml:voice ... TEXT
  synth custom   (stub; not implemented yet)

-v / --voice and -P / --voice-path are mutually exclusive.`,
		Args:                  cobra.ArbitraryArgs,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSynth(cmd, app, opts, args)
		},
	}

	cmd.Flags().StringP("voice", "v", "", "Voice ref character:voice")
	cmd.Flags().StringP("voice-path", "P", "", "Character YAML override path.yaml:voice")
	cmd.Flags().IntP("version", "V", 0, "Engine version 1|2|10 (default: first available)")
	cmd.Flags().Float64P("speed", "s", 1, "Speed factor multiplied with voice speed")
	cmd.Flags().StringP("output", "o", "", "Output WAV path")
	cmd.Flags().BoolP("preview", "p", false, "Play synthesized audio")
	cmd.Flags().Bool("reuse-cache", true, "Reuse cached WAV when task hash matches (overrides synth.reuse_cache)")

	cmd.AddCommand(newVoiceCmd(app, opts))
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

func runSynth(cmd *cobra.Command, app *api.App, opts Options, args []string) error {
	voiceRef, _ := cmd.Flags().GetString("voice")
	voicePath, _ := cmd.Flags().GetString("voice-path")
	version, _ := cmd.Flags().GetInt("version")
	speed, _ := cmd.Flags().GetFloat64("speed")
	output, _ := cmd.Flags().GetString("output")
	preview, _ := cmd.Flags().GetBool("preview")

	if strings.TrimSpace(voiceRef) != "" && strings.TrimSpace(voicePath) != "" {
		return fmt.Errorf("--voice/-v and --voice-path/-P are mutually exclusive")
	}

	text := strings.TrimSpace(strings.Join(args, " "))
	if err := prompt.EnsureString("text", &text); err != nil {
		return err
	}
	if text == "" {
		return fmt.Errorf("text is required")
	}
	if strings.TrimSpace(voicePath) == "" {
		if err := ensureVoiceRef(app, opts.Interactive, &voiceRef); err != nil {
			return err
		}
	}

	req := api.SynthRequest{
		VoiceRef:      voiceRef,
		VoicePathSpec: voicePath,
		Text:          text,
		Version:       version,
		SpeedFactor:   speed,
		Output:        output,
		Preview:       preview,
	}
	if cmd.Flags().Changed("reuse-cache") {
		v, _ := cmd.Flags().GetBool("reuse-cache")
		req.ReuseCache = &v
	}

	_, _, err := runRequest(cmd, app, req, false)
	return err
}
