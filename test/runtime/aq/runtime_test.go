package aq_test

import (
	"context"
	"testing"

	"yukumo-script/internal/domain"
	"yukumo-script/internal/runtime/aq"
)

func TestPresetNamesAndHasPreset(t *testing.T) {
	ctx := context.Background()
	rt, err := aq.New(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = rt.Close(ctx) })

	for _, ver := range []int32{domain.Version1, domain.Version2, domain.Version10} {
		names, err := rt.PresetNames(ctx, ver)
		if err != nil {
			t.Fatalf("version %d names: %v", ver, err)
		}
		if len(names) == 0 {
			t.Fatalf("version %d: empty names", ver)
		}
		ok, err := rt.HasPreset(ctx, ver, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatalf("version %d preset 0 should be valid", ver)
		}
		ok, err = rt.HasPreset(ctx, ver, int32(len(names)+10))
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Fatalf("version %d out-of-range preset should be invalid", ver)
		}
		t.Logf("v%d: %v", ver, names)
	}
}
