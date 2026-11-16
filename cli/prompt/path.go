package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PathOpts struct {
	Title       string
	Description string
	DirAllowed  bool
	FileAllowed bool
	StartDir    string
	ShowHidden  bool
}

const (
	pathActType = iota
	pathActUse
	pathActUp
	pathActEntry
)

type pathItem struct {
	kind  int
	label string
	name  string
	isDir bool
}

func Path(opts PathOpts, value *string) error {
	if value == nil {
		return fmt.Errorf("path value is nil")
	}

	dir := opts.StartDir
	if dir == "" && *value != "" {
		dir = pathStart(*value)
	}
	if dir == "" {
		dir, _ = os.Getwd()
	}
	dir = ExpandPath(dir)

	for {
		items, err := buildPathItems(dir, opts)
		if err != nil {
			return err
		}
		labels := make([]string, len(items))
		for i, it := range items {
			labels[i] = it.label
		}

		idx, err := Select(opts.Title, pathDescription(opts.Description, dir), labels, 0)
		if err != nil {
			return err
		}
		if idx < 0 || idx >= len(items) {
			return fmt.Errorf("invalid selection")
		}

		it := items[idx]
		switch it.kind {
		case pathActType:
			typed := *value
			if typed == "" {
				typed = dir
			}
			if err := String("Enter path", &typed); err != nil {
				return err
			}
			*value = ExpandPath(typed)
			return nil
		case pathActUse:
			*value = ExpandPath(dir)
			return nil
		case pathActUp:
			parent := filepath.Dir(dir)
			if parent != dir {
				dir = parent
			}
		case pathActEntry:
			next := ExpandPath(filepath.Join(dir, it.name))
			if it.isDir {
				dir = next
				continue
			}
			*value = next
			return nil
		}
	}
}

func Dir(title string, value *string) error {
	start := ""
	if value != nil && *value != "" {
		start = pathStart(*value)
	}
	return Path(PathOpts{
		Title:       title,
		Description: "↑↓ move · enter confirm",
		DirAllowed:  true,
		FileAllowed: false,
		StartDir:    start,
	}, value)
}

func File(title string, value *string) error {
	return Path(PathOpts{
		Title:       title,
		Description: "↑↓ move · enter confirm",
		DirAllowed:  false,
		FileAllowed: true,
	}, value)
}

func buildPathItems(dir string, opts PathOpts) ([]pathItem, error) {
	items := []pathItem{
		{kind: pathActUp, label: LabelParent},
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		entries = nil
	}
	for _, e := range entries {
		name := e.Name()
		if !opts.ShowHidden && strings.HasPrefix(name, ".") {
			continue
		}
		if e.IsDir() {
			items = append(items, pathItem{
				kind:  pathActEntry,
				label: name + "/",
				name:  name,
				isDir: true,
			})
			continue
		}
		if opts.FileAllowed {
			items = append(items, pathItem{
				kind:  pathActEntry,
				label: name,
				name:  name,
				isDir: false,
			})
		}
	}

	if opts.DirAllowed {
		items = append(items, pathItem{kind: pathActUse, label: "✓  Use this directory"})
	}
	items = append(items, pathItem{kind: pathActType, label: "✎  Type path…"})
	return items, nil
}

func pathDescription(desc, dir string) string {
	if desc == "" {
		return dir
	}
	return desc + "\n" + dir
}

func pathStart(current string) string {
	current = ExpandPath(current)
	info, err := os.Stat(current)
	if err == nil {
		if info.IsDir() {
			return current
		}
		return filepath.Dir(current)
	}
	parent := filepath.Dir(current)
	if parent != "" && parent != current {
		if _, err := os.Stat(parent); err == nil {
			return parent
		}
	}
	return ""
}

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
