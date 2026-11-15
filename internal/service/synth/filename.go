package synth

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

// FormatOutputName expands tokens in format (default c-v-ymd).
// Braced tokens: {c} {v} {ymd} {ymdhms}. Compact form also supports c, v, ymd, ymdhms.
func FormatOutputName(format, character, voice string) string {
	format = strings.TrimSpace(format)
	if format == "" {
		format = "c-v-ymd"
	}
	now := time.Now()
	repl := map[string]string{
		"c":      sanitizeName(character),
		"v":      sanitizeName(voice),
		"ymd":    now.Format("20060102"),
		"ymdhms": now.Format("20060102150405"),
	}
	if strings.Contains(format, "{") {
		name := format
		for _, key := range []string{"ymdhms", "ymd", "c", "v"} {
			name = strings.ReplaceAll(name, "{"+key+"}", repl[key])
		}
		if name == "" {
			return "synth"
		}
		return name
	}
	var b strings.Builder
	for i := 0; i < len(format); {
		switch {
		case strings.HasPrefix(format[i:], "ymdhms"):
			b.WriteString(repl["ymdhms"])
			i += 6
		case strings.HasPrefix(format[i:], "ymd"):
			b.WriteString(repl["ymd"])
			i += 3
		case format[i] == 'c':
			b.WriteString(repl["c"])
			i++
		case format[i] == 'v':
			b.WriteString(repl["v"])
			i++
		default:
			b.WriteByte(format[i])
			i++
		}
	}
	name := b.String()
	if name == "" {
		name = "synth"
	}
	return name
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
