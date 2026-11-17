package character_test

import (
	"errors"
	"testing"

	"yukumo-script/core/internal/domain"
	svccharacter "yukumo-script/core/internal/service/character"
	"yukumo-script/core/internal/store"
)

func TestCRUD(t *testing.T) {
	mem := store.NewMemory()
	svc := svccharacter.New(nil, mem)

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

func TestAddRequiresName(t *testing.T) {
	svc := svccharacter.New(nil, store.NewMemory())
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
