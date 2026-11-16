package schema_test

import (
	"strings"
	"testing"

	"yukumo-script/cli/schema"
)

func TestNormalizeIntClampsWithWarning(t *testing.T) {
	value, warning, err := schema.Normalize(schema.Field{Kind: schema.KindInt, Min: 1, Max: 10}, "20")
	if err != nil {
		t.Fatal(err)
	}
	if value != "10" {
		t.Fatalf("value = %q, want 10", value)
	}
	if !strings.Contains(warning, "clamped") {
		t.Fatalf("warning = %q, want clamped", warning)
	}
}

func TestNormalizeBool(t *testing.T) {
	value, warning, err := schema.Normalize(schema.Field{Kind: schema.KindBool}, "yes")
	if err != nil {
		t.Fatal(err)
	}
	if value != "true" || warning != "" {
		t.Fatalf("value=%q warning=%q", value, warning)
	}
}

func TestLookupUnknown(t *testing.T) {
	_, err := schema.Lookup(schema.ConfigFields(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
}
