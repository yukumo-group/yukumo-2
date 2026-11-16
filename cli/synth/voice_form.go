package synth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"yukumo-script/cli/prompt"
	"yukumo-script/core/api"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var dimForm = lipgloss.NewStyle().Faint(true)

type voiceDraft struct {
	voiceRef    string
	version     int // 0 = auto
	speedFactor float64
	output      string
	preview     bool
	reuseCache  bool
	text        string
	status      string
}

func newVoiceCmd(app *api.App, opts Options) *cobra.Command {
	return &cobra.Command{
		Use:   "voice",
		Short: "Synthesize from a character voice",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !opts.Interactive {
				return fmt.Errorf("interactive only; use: synth -v character:voice TEXT")
			}
			return runInteractiveVoice(cmd, app)
		},
	}
}

func runInteractiveVoice(cmd *cobra.Command, app *api.App) error {
	d := voiceDraft{speedFactor: 1, reuseCache: true}
	return prompt.With("voice", func() error {
		focus := 0
		for {
			labels := []string{
				prompt.LabelBack,
				fieldLabel("voice", voiceRefDisplay(d.voiceRef)),
				fieldLabel("version", versionDisplay(d.version)),
				fieldLabel("speed", formatFactor(d.speedFactor)),
				fieldLabel("output", outputDisplay(d.output)),
				fieldLabel("preview", strconv.FormatBool(d.preview)),
				fieldLabel("reuse_cache", strconv.FormatBool(d.reuseCache)),
				fieldLabel("text", textDisplay(d.text)),
				"▶  synthesize",
			}
			if d.status != "" {
				labels = append(labels, fieldLabel("last", d.status))
			}
			idx, err := prompt.Select("", "↑↓ move · enter", labels, focus)
			if err != nil {
				if errors.Is(err, prompt.ErrAborted) {
					return nil
				}
				return err
			}
			switch idx {
			case 0:
				return nil
			case 1:
				if err := editVoiceRef(app, &d.voiceRef); err != nil {
					if !errors.Is(err, prompt.ErrAborted) {
						_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
					}
				}
				focus = 1
			case 2:
				if err := editVersion(&d.version); err != nil {
					_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
				}
				focus = 2
			case 3:
				if err := editSpeedFactor(&d.speedFactor); err != nil {
					_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
				}
				focus = 3
			case 4:
				if err := editOutput(&d.output); err != nil {
					if !errors.Is(err, prompt.ErrAborted) {
						_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
					}
				}
				focus = 4
			case 5:
				prompt.ToggleBool(&d.preview)
				focus = 5
			case 6:
				prompt.ToggleBool(&d.reuseCache)
				focus = 6
			case 7:
				if err := prompt.String("text", &d.text); err != nil {
					_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
				}
				focus = 7
			case 8:
				if err := doSynthesize(cmd, app, &d); err != nil {
					_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
					d.status = "error"
				}
				focus = 8
			default:
				focus = 0
			}
		}
	})
}

func doSynthesize(cmd *cobra.Command, app *api.App, d *voiceDraft) error {
	if strings.TrimSpace(d.voiceRef) == "" {
		return fmt.Errorf("voice is required")
	}
	if strings.TrimSpace(d.text) == "" {
		return fmt.Errorf("text is required")
	}
	reuse := d.reuseCache
	res, elapsed, err := runRequest(cmd, app, api.SynthRequest{
		VoiceRef:    d.voiceRef,
		Text:        d.text,
		Version:     d.version,
		SpeedFactor: d.speedFactor,
		Output:      d.output,
		Preview:     d.preview,
		ReuseCache:  &reuse,
	}, true)
	if err != nil {
		return err
	}
	ms := formatElapsed(elapsed)
	switch {
	case res.Wrote && res.FromCache:
		d.status = res.Path + " (cached) " + ms
	case res.Wrote:
		d.status = res.Path + " " + ms
	case res.FromCache:
		d.status = "ok (cached, not written) " + ms
	default:
		d.status = "ok (not written) " + ms
	}
	return nil
}

func editVoiceRef(app *api.App, ref *string) error {
	keys, metas, err := voiceTree(app)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return fmt.Errorf("no character voices")
	}
	return prompt.TreeKeySep("pick", ":", keys, metas, ref)
}

func voiceTree(app *api.App) ([]string, map[string]prompt.LeafMeta, error) {
	list, err := app.ListCharacters()
	if err != nil {
		return nil, nil, err
	}
	keys := make([]string, 0)
	metas := make(map[string]prompt.LeafMeta)
	for _, c := range list {
		vs, err := app.ListVoices(c.Name)
		if err != nil {
			return nil, nil, err
		}
		for _, v := range vs {
			key := c.Name + ":" + v.Name
			keys = append(keys, key)
			metas[key] = prompt.LeafMeta{
				Value:       fmt.Sprintf("%d param(s)", len(v.Params)),
				Description: c.Description,
			}
		}
	}
	return keys, metas, nil
}

func editVersion(version *int) error {
	labels := []string{"auto (1→2→10)", "1", "2", "10"}
	values := []int{0, api.Version1, api.Version2, api.Version10}
	cur := 0
	for i, v := range values {
		if v == *version {
			cur = i
			break
		}
	}
	idx, err := prompt.Select("version", "↑↓ move · enter", labels, cur)
	if err != nil {
		return err
	}
	*version = values[idx]
	return nil
}

func editSpeedFactor(factor *float64) error {
	raw := formatFactor(*factor)
	if err := prompt.String("speed factor", &raw); err != nil {
		return err
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil || v <= 0 {
		return fmt.Errorf("invalid speed factor %q", raw)
	}
	*factor = v
	return nil
}

func editOutput(output *string) error {
	labels := []string{
		prompt.LabelBack,
		"use default format path",
		"browse file…",
	}
	if *output != "" {
		labels = append(labels, "clear  —  "+*output)
	}
	idx, err := prompt.Select("output", "↑↓ move · enter", labels, 0)
	if err != nil {
		return err
	}
	switch idx {
	case 0:
		return prompt.ErrAborted
	case 1:
		*output = ""
		return nil
	case 2:
		v := *output
		if err := prompt.File("output", &v); err != nil {
			return err
		}
		*output = v
		return nil
	case 3:
		*output = ""
		return nil
	}
	return nil
}

func fieldLabel(name, value string) string {
	if value == "" {
		value = "(empty)"
	}
	return name + "\n  " + dimForm.Render(value)
}

func voiceRefDisplay(ref string) string {
	if ref == "" {
		return "(pick character:voice)"
	}
	return ref
}

func versionDisplay(v int) string {
	if v == 0 {
		return "auto"
	}
	return strconv.Itoa(v)
}

func outputDisplay(p string) string {
	if p == "" {
		return "(auto from format)"
	}
	return p
}

func textDisplay(t string) string {
	t = strings.TrimSpace(t)
	if t == "" {
		return "(empty)"
	}
	if len([]rune(t)) > 40 {
		r := []rune(t)
		return string(r[:40]) + "…"
	}
	return t
}

func formatFactor(f float64) string {
	if f == float64(int(f)) {
		return strconv.Itoa(int(f))
	}
	return strconv.FormatFloat(f, 'f', 2, 64)
}
