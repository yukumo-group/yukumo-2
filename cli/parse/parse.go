// Package parse contains pure parsing helpers shared by schema, prompt, and
// shell code. It has no terminal or cobra dependencies.
package parse

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ClampInt clamps v into [min, max] and reports whether it changed.
func ClampInt(v, min, max int) (int, bool) {
	switch {
	case v < min:
		return min, true
	case v > max:
		return max, true
	default:
		return v, false
	}
}

// ParseRangedInt parses raw, clamps it into [min, max], and reports whether it
// was clamped.
func ParseRangedInt(raw string, min, max int) (int, bool, error) {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false, fmt.Errorf("invalid number %q", raw)
	}
	clamped, changed := ClampInt(v, min, max)
	return clamped, changed, nil
}

// ParseBool parses common boolean spellings.
func ParseBool(raw string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "true", "1", "yes", "y", "on":
		return true, nil
	case "false", "0", "no", "n", "off":
		return false, nil
	default:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return false, fmt.Errorf("invalid bool %q", raw)
		}
		return v, nil
	}
}

// ExpandPath expands ~ and returns an absolute, cleaned path.
func ExpandPath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return p
	}
	if p == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
		return p
	}
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			p = filepath.Join(home, p[2:])
		}
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return filepath.Clean(p)
	}
	return filepath.Clean(abs)
}

// ReadPhont reads a .phont file and rejects empty payloads.
func ReadPhont(path string) ([]byte, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("phont path is required")
	}
	if !strings.EqualFold(filepath.Ext(path), ".phont") {
		return nil, fmt.Errorf("phont path must end with .phont")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("phont file is empty")
	}
	return data, nil
}
