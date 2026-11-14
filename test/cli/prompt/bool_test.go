package prompt_test

import (
	"testing"

	"yukumo-script/internal/cli/prompt"
)

func TestParseBool(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"true", true},
		{"yes", true},
		{"1", true},
		{"false", false},
		{"no", false},
		{"0", false},
	}
	for _, tc := range cases {
		got, err := prompt.ParseBool(tc.in)
		if err != nil || got != tc.want {
			t.Fatalf("%q: got %v %v", tc.in, got, err)
		}
	}
	if _, err := prompt.ParseBool("maybe"); err == nil {
		t.Fatal("expected error")
	}
}
