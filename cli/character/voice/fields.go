package voice

import (
	"fmt"
	"io"

	"yukumo-script/cli/schema"
	"yukumo-script/core/api"
)

func applyParamField(p *api.SynthParam, key, raw string, notify io.Writer) error {
	entries := schema.VoiceParamFields(p.Version)
	entry, err := schema.Lookup(entries, key)
	if err != nil {
		return err
	}
	value, warning, err := schema.Normalize(entry.Field, raw)
	if err != nil {
		return err
	}
	if warning != "" && notify != nil {
		_, _ = fmt.Fprint(notify, warning)
	}
	if key == "phont_data" {
		return api.SetPhontPath(p, value)
	}
	if err := entry.Set(p, value); err != nil {
		return err
	}
	p.Normalize()
	return nil
}
