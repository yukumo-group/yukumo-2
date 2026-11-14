package prompt_test

import (
	"testing"

	"yukumo-script/internal/cli/prompt"
)

func TestEnsureTreeKeyValidates(t *testing.T) {
	key := "cli.default_interactive"
	keys := []string{
		"paths.character_path",
		"cli.default_interactive",
		"cli.with_accessibility",
	}
	if err := prompt.EnsureTreeKey("Field", keys, map[string]prompt.LeafMeta{}, &key); err != nil {
		t.Fatal(err)
	}
	key = "nope"
	if err := prompt.EnsureTreeKey("Field", keys, map[string]prompt.LeafMeta{}, &key); err == nil {
		t.Fatal("expected error")
	}
}
