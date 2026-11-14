package nameutil_test

import (
	"testing"

	"yukumo-script/internal/service/nameutil"
)

func TestResolveCopy(t *testing.T) {
	if got := nameutil.ResolveCopy("a", ""); got != "a - copy" {
		t.Fatalf("%q", got)
	}
	if got := nameutil.ResolveCopy("a", "b"); got != "b" {
		t.Fatalf("%q", got)
	}
}
