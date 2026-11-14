package character_test

import (
	"context"
	"errors"
	"testing"

	"yukumo-script/internal/domain"
	svccharacter "yukumo-script/internal/service/character"
	"yukumo-script/internal/store"
)

type allowAllPresets struct{}

func (allowAllPresets) HasPreset(context.Context, int32, int32) (bool, error) {
	return true, nil
}
func (allowAllPresets) PresetNames(context.Context, int32) ([]string, error) {
	return []string{"p0"}, nil
}

func TestCRUD(t *testing.T) {
	mem := store.NewMemory()
	svc := svccharacter.NewWithStore(mem, allowAllPresets{})

	if err := svc.Add(domain.Character{Name: "Alice", Description: "hero"}); err != nil {
		t.Fatal(err)
	}
	if err := svc.Add(domain.Character{Name: "Alice"}); !errors.Is(err, domain.ErrExists) {
		t.Fatalf("duplicate: %v", err)
	}
	if err := svc.Add(domain.Character{Name: "a/b"}); err != nil {
		t.Fatal(err)
	}
	if err := svc.Add(domain.Character{Name: "a:b"}); !errors.Is(err, domain.ErrExists) {
		t.Fatalf("sanitized collision: %v", err)
	}

	got, err := svc.Get("Alice")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Alice" || got.Description != "hero" {
		t.Fatalf("get = %+v", got)
	}

	list, err := svc.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0].Name != "a/b" || list[1].Name != "Alice" {
		t.Fatalf("list = %+v", list)
	}

	if err := svc.Delete("a/b"); err != nil {
		t.Fatal(err)
	}

	if err := svc.Update("Alice", domain.Character{Name: "Bob", Description: "renamed"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Get("Alice"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("old name still present: %v", err)
	}
	got, err = svc.Get("Bob")
	if err != nil {
		t.Fatal(err)
	}
	if got.Description != "renamed" {
		t.Fatalf("after rename = %+v", got)
	}

	if err := svc.Add(domain.Character{Name: "Carol"}); err != nil {
		t.Fatal(err)
	}
	list, err = svc.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0].Name != "Bob" || list[1].Name != "Carol" {
		t.Fatalf("sorted list = %+v", list)
	}

	if err := svc.Delete("Bob"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Get("Bob"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("deleted still get: %v", err)
	}
	if err := svc.Delete("Bob"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("delete missing: %v", err)
	}
}

func TestVoiceCRUD(t *testing.T) {
	svc := svccharacter.NewWithStore(store.NewMemory(), allowAllPresets{})
	if err := svc.Add(domain.Character{Name: "Alice"}); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	v := domain.Voice{
		Name: "default",
		Params: []domain.SynthParam{
			{Version: domain.Version1, Preset: 0, Speed: 100, Volume: 999, PhontData: []byte("x")},
		},
	}
	if err := svc.AddVoice(ctx, "Alice", v); err != nil {
		t.Fatal(err)
	}
	got, err := svc.GetVoice("Alice", "default")
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
	if err := svc.UpdateVoice(ctx, "Alice", "default", got); err != nil {
		t.Fatal(err)
	}
	got, err = svc.GetVoice("Alice", "default")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Params) != 2 || got.Params[1].PhontData != nil {
		t.Fatalf("normalize v10 = %+v", got.Params[1])
	}

	if err := svc.AddVoice(ctx, "Alice", domain.Voice{Name: "default"}); !errors.Is(err, domain.ErrExists) {
		t.Fatalf("dup voice: %v", err)
	}
	if err := svc.DeleteVoice("Alice", "default"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.GetVoice("Alice", "default"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("deleted voice: %v", err)
	}
}

func TestAddRequiresName(t *testing.T) {
	svc := svccharacter.NewWithStore(store.NewMemory(), allowAllPresets{})
	if err := svc.Add(domain.Character{}); !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("err=%v", err)
	}
}

func TestSynthParamNormalize(t *testing.T) {
	p := domain.SynthParam{Version: domain.Version2, Speed: 999, Base: 9, PhontData: []byte{1}}
	p.Normalize()
	if p.Speed != 300 || p.Base != 0 || len(p.PhontData) != 1 {
		t.Fatalf("%+v", p)
	}
	p = domain.SynthParam{Version: domain.Version10, Speed: 10, Pitch: 5, PhontData: []byte{1}}
	p.Normalize()
	if p.Speed != 50 || p.Pitch != 20 || p.PhontData != nil {
		t.Fatalf("%+v", p)
	}
}
