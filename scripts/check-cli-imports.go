// Command check-cli-imports verifies the CLI package layering that depguard
// can only check at the direct-import level.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type pkg struct {
	ImportPath string
	Imports    []string
}

func main() {
	cmd := exec.Command("go", "list", "-json", "./...")
	cmd.Dir = "cli"
	cmd.Env = append(os.Environ(), "GOWORK=off")
	out, err := cmd.StdoutPipe()
	if err != nil {
		fail(err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fail(err)
	}

	pkgs := map[string]pkg{}
	dec := json.NewDecoder(out)
	for {
		var p pkg
		if err := dec.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			fail(err)
		}
		pkgs[p.ImportPath] = p
	}
	if err := cmd.Wait(); err != nil {
		fail(err)
	}

	closure := func(path string) map[string]bool {
		seen := map[string]bool{}
		var walk func(string)
		walk = func(p string) {
			if seen[p] {
				return
			}
			seen[p] = true
			for _, imp := range pkgs[p].Imports {
				walk(imp)
			}
		}
		walk(path)
		return seen
	}

	var problems []string
	for path := range pkgs {
		if !strings.HasPrefix(path, "yukumo-script/cli/") || strings.HasPrefix(path, "yukumo-script/cli/cmd/") {
			continue
		}
		if closure(path)["yukumo-script/cli"] {
			problems = append(problems, fmt.Sprintf("%s must not depend on yukumo-script/cli", path))
		}
	}

	rules := []struct {
		pkg  string
		deny []string
	}{
		{"yukumo-script/cli/action", []string{
			"github.com/spf13/cobra", "github.com/charmbracelet/huh", "github.com/charmbracelet/lipgloss",
			"yukumo-script/cli/prompt", "yukumo-script/cli/setting", "yukumo-script/cli/ui",
			"yukumo-script/cli/render", "yukumo-script/cli/play", "yukumo-script/cli/streams",
			"yukumo-script/cli/menu", "yukumo-script/cli/internal",
			"yukumo-script/cli/config", "yukumo-script/cli/character", "yukumo-script/cli/synth",
		}},
		{"yukumo-script/cli/ui", []string{
			"github.com/spf13/cobra", "yukumo-script/cli/config", "yukumo-script/cli/character",
			"yukumo-script/cli/synth", "yukumo-script/cli/internal",
		}},
		{"yukumo-script/cli/render", []string{
			"github.com/spf13/cobra", "github.com/charmbracelet/huh", "github.com/charmbracelet/lipgloss",
			"yukumo-script/cli/prompt", "yukumo-script/cli/setting", "yukumo-script/cli/ui",
			"yukumo-script/cli/action", "yukumo-script/cli/menu", "yukumo-script/cli/streams",
			"yukumo-script/cli/internal",
		}},
		{"yukumo-script/cli/play", []string{
			"github.com/spf13/cobra", "github.com/charmbracelet/huh", "github.com/charmbracelet/lipgloss",
			"yukumo-script/cli/prompt", "yukumo-script/cli/setting", "yukumo-script/cli/ui",
			"yukumo-script/cli/action", "yukumo-script/cli/menu", "yukumo-script/cli/streams",
			"yukumo-script/cli/internal",
		}},
		{"yukumo-script/cli/menu", []string{
			"github.com/spf13/cobra", "github.com/charmbracelet/huh", "github.com/charmbracelet/lipgloss",
			"yukumo-script/cli/prompt", "yukumo-script/cli/setting", "yukumo-script/cli/ui",
			"yukumo-script/cli/action", "yukumo-script/cli/schema", "yukumo-script/cli/internal",
		}},
		{"yukumo-script/cli/prompt", []string{
			"github.com/spf13/cobra", "yukumo-script/core/api", "yukumo-script/cli/schema",
			"yukumo-script/cli/setting", "yukumo-script/cli/action", "yukumo-script/cli/ui",
		}},
		{"yukumo-script/cli/setting", []string{
			"github.com/spf13/cobra", "yukumo-script/cli/action", "yukumo-script/cli/ui",
			"yukumo-script/cli/internal",
		}},
		{"yukumo-script/cli/schema", []string{
			"github.com/spf13/cobra", "github.com/charmbracelet/huh", "github.com/charmbracelet/lipgloss",
			"yukumo-script/cli/prompt", "yukumo-script/cli/setting", "yukumo-script/cli/ui",
			"yukumo-script/cli/action", "yukumo-script/cli/internal",
		}},
	}

	for _, rule := range rules {
		deps := closure(rule.pkg)
		for _, bad := range rule.deny {
			if deps[bad] || hasPrefix(deps, bad) {
				problems = append(problems, fmt.Sprintf("%s must not depend on %s", rule.pkg, bad))
			}
		}
		if strings.HasSuffix(rule.pkg, "/ui") || strings.HasSuffix(rule.pkg, "/action") {
			if deps["yukumo-script/cli"] {
				problems = append(problems, rule.pkg+" must not depend on yukumo-script/cli")
			}
		}
	}

	for _, leaf := range []string{"yukumo-script/cli/parse", "yukumo-script/cli/streams"} {
		for dep := range closure(leaf) {
			if dep != leaf && strings.HasPrefix(dep, "yukumo-script/") {
				problems = append(problems, fmt.Sprintf("%s must stay stdlib-only, found %s", leaf, dep))
			}
		}
	}

	if len(problems) > 0 {
		for _, p := range problems {
			fmt.Fprintln(os.Stderr, p)
		}
		os.Exit(1)
	}
	fmt.Println("cli import graph ok")
}

func hasPrefix(deps map[string]bool, prefix string) bool {
	for dep := range deps {
		if strings.HasPrefix(dep, prefix+"/") {
			return true
		}
	}
	return false
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
