package prompt_test

import (
	"testing"

	"yukumo-script/internal/cli/prompt"
)

func TestNavTitleAndWith(t *testing.T) {
	prompt.Reset()
	if got := prompt.Title("alone"); got != "alone" {
		t.Fatalf("empty crumb Title: got %q", got)
	}

	prompt.Reset("root", "config")
	if got := prompt.Crumb(); got != "root › config" {
		t.Fatalf("Crumb: got %q", got)
	}
	if got := prompt.Title(""); got != "root › config" {
		t.Fatalf("Title empty segment: got %q", got)
	}
	if got := prompt.Title("edit"); got != "root › config › edit" {
		t.Fatalf("Title with segment: got %q", got)
	}

	err := prompt.With("edit", func() error {
		if got := prompt.Crumb(); got != "root › config › edit" {
			t.Fatalf("inside With: got %q", got)
		}
		prompt.Push("paths")
		prompt.Push("character_path")
		if got := prompt.Crumb(); got != "root › config › edit › paths › character_path" {
			t.Fatalf("nested push: got %q", got)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := prompt.Crumb(); got != "root › config" {
		t.Fatalf("after With restore: got %q", got)
	}
	prompt.Reset()
}
