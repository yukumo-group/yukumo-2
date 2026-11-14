package character

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"yukumo-script/internal/domain"
)

func LoadPhont(path string) ([]byte, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("%w: phont path is required", domain.ErrInvalid)
	}
	if !strings.EqualFold(filepath.Ext(path), ".phont") {
		return nil, fmt.Errorf("%w: phont path must end with .phont", domain.ErrInvalid)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("%w: phont file is empty", domain.ErrInvalid)
	}
	return data, nil
}

func SetPhontPath(p *domain.SynthParam, path string) error {
	if p == nil {
		return fmt.Errorf("%w: nil param", domain.ErrInvalid)
	}
	if strings.TrimSpace(path) == "" {
		p.PhontData = nil
		return nil
	}
	data, err := LoadPhont(path)
	if err != nil {
		return err
	}
	p.PhontData = data
	return nil
}
