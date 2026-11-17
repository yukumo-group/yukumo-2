package character

import (
	"fmt"
	"regexp"
	"strings"

	"yukumo-script/core/internal/domain"
)

var illegalFileChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)

func SanitizeBase(name string) (string, error) {
	s := strings.TrimSpace(name)
	s = illegalFileChars.ReplaceAllString(s, "_")
	s = strings.Trim(s, " .")
	if s == "" {
		return "", fmt.Errorf("%w: name has no valid file characters", domain.ErrInvalid)
	}
	return s, nil
}

func FileName(name string) (string, error) {
	base, err := SanitizeBase(name)
	if err != nil {
		return "", err
	}
	return base + ".yaml", nil
}
