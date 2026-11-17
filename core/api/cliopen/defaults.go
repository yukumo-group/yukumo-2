package cliopen

import (
	"fmt"
	"os"
	"path/filepath"

	"yukumo-script/core/api"
	"yukumo-script/core/internal/domain"
)

func defaultCharacterPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}
	return filepath.Join(base, "yukumo-script", "characters"), nil
}

func defaultOutputPath() string {
	return filepath.Join(os.TempDir(), "yukumo-script")
}

func domainDefaults() (domain.Config, error) {
	charPath, err := defaultCharacterPath()
	if err != nil {
		return domain.Config{}, err
	}
	return domain.DefaultConfig(domain.PathConfig{
		CharacterPath: charPath,
		OutputPath:    defaultOutputPath(),
	}), nil
}

func apiDefaults(d domain.Config) api.Config {
	c := api.DefaultConfig()
	c.Paths = api.PathConfig{
		CharacterPath: d.Paths.CharacterPath,
		OutputPath:    d.Paths.OutputPath,
	}
	return c
}

// Defaults returns the OS-backed API defaults for CLI embedding/tests.
func Defaults() (api.Config, error) {
	d, err := domainDefaults()
	if err != nil {
		return api.Config{}, err
	}
	return apiDefaults(d), nil
}
