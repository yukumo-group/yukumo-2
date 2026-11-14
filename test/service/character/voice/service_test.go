package voice_test

import (
	"context"
	"errors"
	"testing"

	"yukumo-script/internal/domain"
	svccharacter "yukumo-script/internal/service/character"
	svcvoice "yukumo-script/internal/service/character/voice"
	"yukumo-script/internal/store"
)

type allowAllPresets struct{}

func (allowAllPresets) HasPreset(context.Context, int32, int32) (bool, error) {
	return true, nil
}
func (allowAllPresets) PresetNames(context.Context, int32) ([]string, error) {
	return []string{"p0"}, nil
}

func TestVoiceCRUD(t *testing.T) {
	chars := svccharacter.NewWithStore(store.NewMemory())
	voices := svcvoice.New(chars, allowAllPresets{})
	if err := chars.Add(domain.Character{Name: "Alice"}); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	v := domain.Voice{
		Name: "default",
		Params: []domain.SynthParam{
			{Version: domain.Version1, Preset: 0, Speed: 100, Volume: 999, PhontData: []byte("x")},
		},
	}
	if err := voices.Add(ctx, "Alice", v); err != nil {
		t.Fatal(err)
	}
	got, err := voices.Get("Alice", "default")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Params) != 1 || got.Params[0].Volume != 0 || got.Params[0].PhontData != nil {
		t.Fatalf("normalize v1 = %+v", got.Params[0])
	}

	got.Params = append(got.Params, domain.SynthParam{
		Version: domain.Version10, Preset: 1, Speed: 100, Volume: 100, Pitch: 100, Accent: 100, Lmd: 100, Fsc: 100,
		PhontData: []byte("nope"),
	})
	if err := voices.Update(ctx, "Alice", "default", got); err != nil {
		t.Fatal(err)
	}
	got, err = voices.Get("Alice", "default")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Params) != 2 || got.Params[1].PhontData != nil {
		t.Fatalf("normalize v10 = %+v", got.Params[1])
	}

	if err := voices.Add(ctx, "Alice", domain.Voice{Name: "default"}); !errors.Is(err, domain.ErrExists) {
		t.Fatalf("dup voice: %v", err)
	}
	if err := voices.Delete("Alice", "default"); err != nil {
		t.Fatal(err)
	}
	if _, err := voices.Get("Alice", "default"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("deleted voice: %v", err)
	}
}
