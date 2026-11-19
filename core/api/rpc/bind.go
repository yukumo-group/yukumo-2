package rpc

import (
	"context"

	"yukumo-script/core/api"
)

func Bind(app *api.App) *Registry {
	if app == nil {
		panic("rpc: app is nil")
	}
	r := NewRegistry()
	r.Register(
		Def("character.list", func(ctx context.Context, _ EmptyReq) (CharacterListRes, error) {
			list, err := app.ListCharacters()
			if err != nil {
				return CharacterListRes{}, err
			}
			return CharacterListRes{Characters: list}, nil
		}),
		Def("character.get", func(ctx context.Context, req CharacterGetReq) (api.Character, error) {
			return app.GetCharacter(req.Name)
		}),
		Def("character.add", func(ctx context.Context, req CharacterAddReq) (EmptyRes, error) {
			return EmptyRes{}, app.AddCharacter(req.Character)
		}),
		Def("character.update", func(ctx context.Context, req CharacterUpdateReq) (EmptyRes, error) {
			return EmptyRes{}, app.UpdateCharacter(req.OldName, req.Character)
		}),
		Def("character.delete", func(ctx context.Context, req CharacterDeleteReq) (EmptyRes, error) {
			return EmptyRes{}, app.DeleteCharacter(req.Name)
		}),

		Def("voice.list", func(ctx context.Context, req VoiceListReq) (VoiceListRes, error) {
			list, err := app.ListVoices(req.Character)
			if err != nil {
				return VoiceListRes{}, err
			}
			return VoiceListRes{Voices: list}, nil
		}),
		Def("voice.get", func(ctx context.Context, req VoiceGetReq) (api.Voice, error) {
			return app.GetVoice(req.Character, req.Name)
		}),
		Def("voice.add", func(ctx context.Context, req VoiceAddReq) (EmptyRes, error) {
			return EmptyRes{}, app.AddVoice(ctx, req.Character, req.Voice)
		}),
		Def("voice.update", func(ctx context.Context, req VoiceUpdateReq) (EmptyRes, error) {
			return EmptyRes{}, app.UpdateVoice(ctx, req.Character, req.OldName, req.Voice)
		}),
		Def("voice.delete", func(ctx context.Context, req VoiceDeleteReq) (EmptyRes, error) {
			return EmptyRes{}, app.DeleteVoice(req.Character, req.Name)
		}),
		Def("voice.duplicate", func(ctx context.Context, req VoiceDuplicateReq) (VoiceDuplicateRes, error) {
			name, err := app.DuplicateVoice(ctx, req.Character, req.SrcName, req.NewName)
			if err != nil {
				return VoiceDuplicateRes{}, err
			}
			return VoiceDuplicateRes{Name: name}, nil
		}),
		Def("voice.presets", func(ctx context.Context, req VoicePresetsReq) (VoicePresetsRes, error) {
			names, err := app.PresetNames(ctx, req.Version)
			if err != nil {
				return VoicePresetsRes{}, err
			}
			return VoicePresetsRes{Names: names}, nil
		}),
		Def("voice.hasPreset", func(ctx context.Context, req VoiceHasPresetReq) (VoiceHasPresetRes, error) {
			ok, err := app.HasPreset(ctx, req.Version, req.Preset)
			if err != nil {
				return VoiceHasPresetRes{}, err
			}
			return VoiceHasPresetRes{OK: ok}, nil
		}),

		Def("config.get", func(ctx context.Context, _ EmptyReq) (api.Config, error) {
			return app.Config()
		}),
		Def("config.save", func(ctx context.Context, req ConfigSaveReq) (EmptyRes, error) {
			return EmptyRes{}, app.SaveConfig(req.Config)
		}),

		Def("synth.render", func(ctx context.Context, req api.SynthRequest) (api.SynthResult, error) {
			req.Preview = false
			return app.Synthesize(ctx, req)
		}),
	)
	return r
}
