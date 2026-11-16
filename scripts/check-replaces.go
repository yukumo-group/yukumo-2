package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type workFile struct {
	Use []struct {
		DiskPath string
	}
}

type modFile struct {
	Module struct {
		Path string
	}
	Replace []struct {
		Old struct {
			Path    string
			Version string
		}
		New struct {
			Path    string
			Version string
		}
	}
}

type target struct {
	modulePath string
	dir        string
	from       string
}

func main() {
	root, err := findRoot()
	if err != nil {
		fail(err)
	}
	modules, err := workModules(root)
	if err != nil {
		fail(err)
	}
	byModule := map[string]target{}
	for _, dir := range modules {
		mf, err := readMod(dir)
		if err != nil {
			fail(err)
		}
		for _, r := range mf.Replace {
			if r.New.Version != "" {
				continue
			}
			old := strings.TrimSpace(r.Old.Path)
			if old == "" {
				fail(fmt.Errorf("%s: empty replace source", dir))
			}
			targetDir, err := resolve(dir, r.New.Path)
			if err != nil {
				fail(fmt.Errorf("%s: replace %s: %w", dir, old, err))
			}
			targetMod, err := readMod(targetDir)
			if err != nil {
				fail(fmt.Errorf("replace %s -> %s: %w", old, targetDir, err))
			}
			if targetMod.Module.Path != old {
				fail(fmt.Errorf("replace %s points to %s, whose module is %q", old, targetDir, targetMod.Module.Path))
			}
			got := target{modulePath: old, dir: targetDir, from: dir}
			if prev, ok := byModule[old]; ok && prev.dir != got.dir {
				fail(fmt.Errorf("replace %s is inconsistent:\n  %s -> %s\n  %s -> %s", old, prev.from, prev.dir, got.from, got.dir))
			}
			byModule[old] = got
		}
	}
	if len(byModule) == 0 {
		fmt.Println("no local replace directives")
		return
	}
	keys := make([]string, 0, len(byModule))
	for k := range byModule {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		rel, _ := filepath.Rel(root, byModule[k].dir)
		fmt.Printf("%-28s -> %s\n", k, rel)
	}
}

func findRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.work not found")
		}
		dir = parent
	}
}

func workModules(root string) ([]string, error) {
	out, err := run(root, "auto", "go", "work", "edit", "-json")
	if err != nil {
		return nil, err
	}
	var wf workFile
	if err := json.Unmarshal(out, &wf); err != nil {
		return nil, fmt.Errorf("parse go.work: %w", err)
	}
	dirs := make([]string, 0, len(wf.Use))
	for _, u := range wf.Use {
		dir, err := resolve(root, u.DiskPath)
		if err != nil {
			return nil, err
		}
		dirs = append(dirs, dir)
	}
	return dirs, nil
}

func readMod(dir string) (modFile, error) {
	var mf modFile
	out, err := run(dir, "off", "go", "mod", "edit", "-json")
	if err != nil {
		return mf, err
	}
	if err := json.Unmarshal(out, &mf); err != nil {
		return mf, fmt.Errorf("parse %s/go.mod: %w", dir, err)
	}
	if mf.Module.Path == "" {
		return mf, fmt.Errorf("%s/go.mod: empty module path", dir)
	}
	return mf, nil
}

func run(dir, gowork string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = envWithGOWORK(gowork)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

func envWithGOWORK(value string) []string {
	env := os.Environ()
	filtered := env[:0]
	for _, e := range env {
		if !strings.HasPrefix(e, "GOWORK=") {
			filtered = append(filtered, e)
		}
	}
	return append(filtered, "GOWORK="+value)
}

func resolve(base, path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(base, path)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if real, err := filepath.EvalSymlinks(abs); err == nil {
		return real, nil
	}
	return abs, nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "check-replaces:", err)
	os.Exit(1)
}
