package action

import (
	"context"
	"strconv"
	"strings"

	"yukumo-script/cli/parse"
	"yukumo-script/cli/schema"
	"yukumo-script/core/api"
)

func ListVoices(_ context.Context, app *api.App, charName string) ([]api.Voice, error) {
	return app.ListVoices(charName)
}

func GetVoice(_ context.Context, app *api.App, charName, name string) (api.Voice, error) {
	return app.GetVoice(charName, name)
}

func AddVoice(ctx context.Context, app *api.App, charName string, v api.Voice) (api.Voice, error) {
	if err := app.AddVoice(ctx, charName, v); err != nil {
		return api.Voice{}, err
	}
	return v, nil
}

func SaveVoice(ctx context.Context, app *api.App, charName, oldName string, v api.Voice) (api.Voice, error) {
	if err := app.UpdateVoice(ctx, charName, oldName, v); err != nil {
		return api.Voice{}, err
	}
	return v, nil
}

func ApplyVoiceName(ctx context.Context, app *api.App, charName, oldName, newName string) (api.Voice, EditResult, error) {
	v, err := app.GetVoice(charName, oldName)
	if err != nil {
		return api.Voice{}, EditResult{}, err
	}
	v.Name = newName
	if err := app.UpdateVoice(ctx, charName, oldName, v); err != nil {
		return api.Voice{}, EditResult{}, err
	}
	return v, EditResult{Value: v.Name}, nil
}

func ApplyVoiceSpeed(ctx context.Context, app *api.App, charName, voiceName, raw string, normalize bool) (api.Voice, EditResult, error) {
	v, err := app.GetVoice(charName, voiceName)
	if err != nil {
		return api.Voice{}, EditResult{}, err
	}
	value := raw
	var warning string
	if normalize {
		value, warning, err = schema.Normalize(schema.Field{
			Key: "speed", Title: "Speed", Description: "Speech rate %", Kind: schema.KindInt, Min: 50, Max: 300,
		}, raw)
		if err != nil {
			return api.Voice{}, EditResult{}, err
		}
	}
	n, _, err := parse.ParseRangedInt(value, 50, 300)
	if err != nil {
		return api.Voice{}, EditResult{}, err
	}
	v.Speed = n
	if err := app.UpdateVoice(ctx, charName, voiceName, v); err != nil {
		return api.Voice{}, EditResult{}, err
	}
	return v, EditResult{Value: strconv.Itoa(v.Speed), Warning: warning}, nil
}

func ApplyVoiceParam(
	ctx context.Context,
	app *api.App,
	charName, voiceName string,
	version int,
	field, raw string,
	normalize bool,
) (api.Voice, EditResult, error) {
	v, err := app.GetVoice(charName, voiceName)
	if err != nil {
		return api.Voice{}, EditResult{}, err
	}
	p, idx := paramByVersion(&v, version)
	if idx < 0 {
		draft := api.DefaultSynthParam(version)
		p = &draft
	}
	entry, err := schema.Lookup(schema.VoiceParamFields(version), field)
	if err != nil {
		return api.Voice{}, EditResult{}, err
	}
	value := raw
	var warning string
	if normalize {
		value, warning, err = schema.Normalize(entry.Field, raw)
		if err != nil {
			return api.Voice{}, EditResult{}, err
		}
	}
	if field == "phont_data" {
		if strings.TrimSpace(value) == "" {
			if err := api.SetPhontData(p, nil); err != nil {
				return api.Voice{}, EditResult{}, err
			}
		} else {
			data, err := parse.ReadPhont(value)
			if err != nil {
				return api.Voice{}, EditResult{}, err
			}
			if err := api.SetPhontData(p, data); err != nil {
				return api.Voice{}, EditResult{}, err
			}
		}
	} else {
		if err := entry.Set(p, value); err != nil {
			return api.Voice{}, EditResult{}, err
		}
		p.Normalize()
	}
	if idx < 0 {
		v.Params = append(v.Params, *p)
	} else {
		v.Params[idx] = *p
	}
	if err := app.UpdateVoice(ctx, charName, voiceName, v); err != nil {
		return api.Voice{}, EditResult{}, err
	}
	return v, EditResult{Value: entry.Get(*p), Warning: warning}, nil
}

func DeleteVoice(_ context.Context, app *api.App, charName, name string) error {
	return app.DeleteVoice(charName, name)
}

func DuplicateVoice(ctx context.Context, app *api.App, charName, src, newName string) (string, error) {
	return app.DuplicateVoice(ctx, charName, src, newName)
}

func paramByVersion(v *api.Voice, version int) (*api.SynthParam, int) {
	for i := range v.Params {
		if v.Params[i].Version == version {
			return &v.Params[i], i
		}
	}
	return nil, -1
}
