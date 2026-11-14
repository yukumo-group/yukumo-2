package prompt

import "strings"

var crumbs []string

func Reset(parts ...string) {
	crumbs = append([]string(nil), parts...)
}

func Push(part string) {
	if part == "" {
		return
	}
	crumbs = append(crumbs, part)
}

func Pop() {
	if len(crumbs) == 0 {
		return
	}
	crumbs = crumbs[:len(crumbs)-1]
}

func Replace(part string) {
	if part == "" {
		return
	}
	if len(crumbs) == 0 {
		crumbs = []string{part}
		return
	}
	crumbs[len(crumbs)-1] = part
}

func With(part string, fn func() error) error {
	before := len(crumbs)
	Push(part)
	defer func() {
		if before < len(crumbs) {
			crumbs = crumbs[:before]
		} else if before == 0 {
			crumbs = nil
		}
	}()
	return fn()
}

func Crumb() string {
	return strings.Join(crumbs, " › ")
}

func Title(segment string) string {
	base := Crumb()
	switch {
	case base == "":
		return segment
	case segment == "":
		return base
	default:
		return base + " › " + segment
	}
}

func syncSuffix(baseLen int, parts []string) {
	if baseLen < 0 {
		baseLen = 0
	}
	if baseLen > len(crumbs) {
		baseLen = len(crumbs)
	}
	crumbs = append(crumbs[:baseLen:baseLen], parts...)
}
