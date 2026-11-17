package synth

import (
	"fmt"
	"os"
	"strings"
)

// parseVoicePathSpec splits "/path/to/character.yaml:voice" on the last colon.
func parseVoicePathSpec(spec string) (yamlPath, voiceName string, err error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return "", "", fmt.Errorf("voice-path is required (file.yaml:voice)")
	}
	i := strings.LastIndexByte(spec, ':')
	if i <= 0 || i == len(spec)-1 {
		return "", "", fmt.Errorf("voice-path must be file.yaml:voice")
	}
	return spec[:i], spec[i+1:], nil
}

func loadVoiceFile(yamlPath, voiceName string) (data []byte, name, base string, err error) {
	data, err = os.ReadFile(yamlPath)
	if err != nil {
		return nil, "", "", err
	}
	return data, strings.TrimSpace(voiceName), baseNameNoExt(yamlPath), nil
}

func baseNameNoExt(path string) string {
	base := path
	if i := strings.LastIndexAny(path, `/\`); i >= 0 {
		base = path[i+1:]
	}
	if i := strings.LastIndexByte(base, '.'); i > 0 {
		return base[:i]
	}
	return base
}
