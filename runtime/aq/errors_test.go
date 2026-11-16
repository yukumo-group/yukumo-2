package aq

import (
	"context"
	"testing"
)

func TestSynthesizeErrorMessageHostCodes(t *testing.T) {
	ctx := context.Background()
	cases := map[int32]string{
		aqErrArg:     "AQ_ERR_ARG",
		aqErrVersion: "AQ_ERR_VERSION",
		aqErrOpen:    "AQ_ERR_OPEN",
	}
	for code, want := range cases {
		got := synthesizeErrorMessage(ctx, nil, nil, 1, code)
		if got != want {
			t.Fatalf("code %d: got %q want %q", code, got, want)
		}
	}
	if got := synthesizeErrorMessage(ctx, nil, nil, 1, -999); got != "unknown error" {
		t.Fatalf("unknown: %q", got)
	}
}
