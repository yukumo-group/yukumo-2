package rpc

import "yukumo-script/core/api"

type EmptyReq struct{}

type EmptyRes struct{}

type CharacterListRes struct {
	Characters []api.Character `json:"characters" jsonschema:"required"`
}

type CharacterGetReq struct {
	Name string `json:"name" jsonschema:"required,minLength=1"`
}

type CharacterAddReq struct {
	Character api.Character `json:"character" jsonschema:"required"`
}

type CharacterUpdateReq struct {
	OldName   string        `json:"oldName" jsonschema:"required,minLength=1"`
	Character api.Character `json:"character" jsonschema:"required"`
}

type CharacterDeleteReq struct {
	Name string `json:"name" jsonschema:"required,minLength=1"`
}

type VoiceListReq struct {
	Character string `json:"character" jsonschema:"required,minLength=1"`
}

type VoiceListRes struct {
	Voices []api.Voice `json:"voices" jsonschema:"required"`
}

type VoiceGetReq struct {
	Character string `json:"character" jsonschema:"required,minLength=1"`
	Name      string `json:"name" jsonschema:"required,minLength=1"`
}

type VoiceAddReq struct {
	Character string    `json:"character" jsonschema:"required,minLength=1"`
	Voice     api.Voice `json:"voice" jsonschema:"required"`
}

type VoiceUpdateReq struct {
	Character string    `json:"character" jsonschema:"required,minLength=1"`
	OldName   string    `json:"oldName" jsonschema:"required,minLength=1"`
	Voice     api.Voice `json:"voice" jsonschema:"required"`
}

type VoiceDeleteReq struct {
	Character string `json:"character" jsonschema:"required,minLength=1"`
	Name      string `json:"name" jsonschema:"required,minLength=1"`
}

type VoiceDuplicateReq struct {
	Character string `json:"character" jsonschema:"required,minLength=1"`
	SrcName   string `json:"srcName" jsonschema:"required,minLength=1"`
	NewName   string `json:"newName" jsonschema:"required,minLength=1"`
}

type VoiceDuplicateRes struct {
	Name string `json:"name" jsonschema:"required"`
}

type VoicePresetsReq struct {
	Version int `json:"version" jsonschema:"required"`
}

type VoicePresetsRes struct {
	Names []string `json:"names" jsonschema:"required"`
}

type VoiceHasPresetReq struct {
	Version int `json:"version" jsonschema:"required"`
	Preset  int `json:"preset" jsonschema:"required"`
}

type VoiceHasPresetRes struct {
	OK bool `json:"ok"`
}

type ConfigSaveReq struct {
	Config api.Config `json:"config" jsonschema:"required"`
}
