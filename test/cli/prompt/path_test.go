package prompt_test

import (
	"os"
	"path/filepath"
	"testing"

	"yukumo-script/internal/cli/prompt"
)

func TestExpandPathHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	got := prompt.ExpandPath("~/foo")
	want := filepath.Join(home, "foo")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestExpandPathRelative(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for _, in := range []string{".", "./", "./."} {
		got := prompt.ExpandPath(in)
		if got != cwd {
			t.Fatalf("%q: got %q want %q", in, got, cwd)
		}
		if !filepath.IsAbs(got) {
			t.Fatalf("%q: not absolute: %q", in, got)
		}
	}
}
