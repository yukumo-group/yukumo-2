package prompt_test

import (
	"testing"

	"yukumo-script/internal/cli/prompt"
)

func TestClampInt(t *testing.T) {
	v, changed := prompt.ClampInt(5, 1, 10)
	if v != 5 || changed {
		t.Fatalf("got %d %v", v, changed)
	}
	v, changed = prompt.ClampInt(-1, 0, 10)
	if v != 0 || !changed {
		t.Fatalf("got %d %v", v, changed)
	}
	v, changed = prompt.ClampInt(99, 0, 10)
	if v != 10 || !changed {
		t.Fatalf("got %d %v", v, changed)
	}
}

func TestParseRangedInt(t *testing.T) {
	v, clamped, err := prompt.ParseRangedInt("150", 0, 100)
	if err != nil || v != 100 || !clamped {
		t.Fatalf("got %d %v %v", v, clamped, err)
	}
}
