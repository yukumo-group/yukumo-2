package synth

import (
	"fmt"
	"strings"

	"yukumo-script/core/internal/domain"

	"gopkg.in/yaml.v3"
)

// ParseVoiceRef splits "character:voice" on the first colon.
func ParseVoiceRef(ref string) (character, voice string, err error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return "", "", fmt.Errorf("%w: voice is required (character:voice)", domain.ErrInvalid)
	}
	i := strings.IndexByte(ref, ':')
	if i <= 0 || i == len(ref)-1 {
		return "", "", fmt.Errorf("%w: voice must be character:voice", domain.ErrInvalid)
	}
	return ref[:i], ref[i+1:], nil
}

// LoadVoiceBytes parses a Character YAML and returns the named voice plus raw
// file bytes (for hashing).
func LoadVoiceBytes(data []byte, voiceName, base string) (domain.Character, domain.Voice, []byte, error) {
	voiceName = strings.TrimSpace(voiceName)
	if voiceName == "" {
		return domain.Character{}, domain.Voice{}, nil, fmt.Errorf("%w: voice name is required", domain.ErrInvalid)
	}
	if len(data) == 0 {
		return domain.Character{}, domain.Voice{}, nil, fmt.Errorf("%w: character file is empty", domain.ErrInvalid)
	}
	var c domain.Character
	if err := yaml.Unmarshal(data, &c); err != nil {
		return domain.Character{}, domain.Voice{}, nil, err
	}
	if c.Name == "" {
		c.Name = strings.TrimSpace(base)
		if c.Name == "" {
			c.Name = "character"
		}
	}
	for _, v := range c.Voices {
		if v.Name == voiceName {
			return c, v, data, nil
		}
	}
	return domain.Character{}, domain.Voice{}, nil, fmt.Errorf("%w: voice %q", domain.ErrNotFound, voiceName)
}

// SelectParam picks a SynthParam: explicit version, else first of 1→2→10 present.
func SelectParam(v domain.Voice, version int) (domain.SynthParam, error) {
	if len(v.Params) == 0 {
		return domain.SynthParam{}, fmt.Errorf("%w: voice %q params empty", domain.ErrInvalid, v.Name)
	}
	byVer := map[int]domain.SynthParam{}
	for _, p := range v.Params {
		byVer[p.Version] = p
	}
	if version != 0 {
		p, ok := byVer[version]
		if !ok {
			return domain.SynthParam{}, fmt.Errorf("%w: voice %q definition invalid (no version %d)", domain.ErrInvalid, v.Name, version)
		}
		return p, nil
	}
	for _, ver := range []int{domain.Version1, domain.Version2, domain.Version10} {
		if p, ok := byVer[ver]; ok {
			return p, nil
		}
	}
	return domain.SynthParam{}, fmt.Errorf("%w: voice %q definition invalid", domain.ErrInvalid, v.Name)
}
