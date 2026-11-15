package synth

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"yukumo-script/internal/domain"
)

var bracedToken = regexp.MustCompile(`\{([^{}]+)\}`)

// FormatOutputName expands braced tokens only (no compact c-v-ymd form).
//
// Named: {character} / {c}, {voice} / {v}
// Time runs (letter length = zero-padded width):
//
//	y… year, M… month, d… day, h|H… hour (24h), m… minute, s… second
//
// {ms} = milliseconds (default 3 digits; longer runs of "ms"+"s" not used — use {ms}).
func FormatOutputName(format, character, voice string) string {
	format = strings.TrimSpace(format)
	if format == "" {
		format = domain.DefaultOutputFormat
	}
	now := time.Now()
	name := bracedToken.ReplaceAllStringFunc(format, func(tok string) string {
		inner := tok[1 : len(tok)-1]
		switch inner {
		case "character", "c":
			return sanitizeName(character)
		case "voice", "v":
			return sanitizeName(voice)
		case "ms":
			return fmt.Sprintf("%03d", now.Nanosecond()/1e6)
		}
		if n, ok := runLen(inner, 'y'); ok {
			return formatYear(now.Year(), n)
		}
		if n, ok := runLen(inner, 'M'); ok {
			return padInt(int(now.Month()), n)
		}
		if n, ok := runLen(inner, 'd'); ok {
			return padInt(now.Day(), n)
		}
		if n, ok := runLen(inner, 'h'); ok {
			return padInt(now.Hour(), n)
		}
		if n, ok := runLen(inner, 'H'); ok {
			return padInt(now.Hour(), n)
		}
		if n, ok := runLen(inner, 'm'); ok {
			return padInt(now.Minute(), n)
		}
		if n, ok := runLen(inner, 's'); ok {
			return padInt(now.Second(), n)
		}
		return tok
	})
	if name == "" {
		return "synth"
	}
	return name
}

func runLen(s string, r byte) (int, bool) {
	if s == "" {
		return 0, false
	}
	for i := 0; i < len(s); i++ {
		if s[i] != r {
			return 0, false
		}
	}
	return len(s), true
}

func formatYear(year, width int) string {
	if width <= 0 {
		width = 4
	}
	s := fmt.Sprintf("%0*d", width, year)
	if len(s) > width {
		return s[len(s)-width:]
	}
	return s
}

func padInt(v, width int) string {
	if width <= 0 {
		width = 1
	}
	return fmt.Sprintf("%0*d", width, v)
}

func sanitizeName(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "unknown"
	}
	var b strings.Builder
	for _, r := range s {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			b.WriteRune(r)
		case r == '-' || r == '_' || r == '.':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	out := b.String()
	if out == "" {
		return "unknown"
	}
	return out
}

// UniquePath returns dir/base.ext, or dir/base-N.ext when the file exists.
func UniquePath(dir, base, ext string) (string, error) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	candidate := filepath.Join(dir, base+ext)
	if _, err := os.Stat(candidate); err != nil {
		if os.IsNotExist(err) {
			return candidate, nil
		}
		return "", err
	}
	for i := 1; i < 10000; i++ {
		candidate = filepath.Join(dir, fmt.Sprintf("%s-%d%s", base, i, ext))
		if _, err := os.Stat(candidate); err != nil {
			if os.IsNotExist(err) {
				return candidate, nil
			}
			return "", err
		}
	}
	return "", fmt.Errorf("could not allocate unique filename for %s", base)
}
