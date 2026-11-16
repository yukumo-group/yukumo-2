package voice_test

import (
	"context"
	"testing"

	"yukumo-script/core/internal/domain"
	svccharacter "yukumo-script/core/internal/service/character"
	svcvoice "yukumo-script/core/internal/service/character/voice"
	"yukumo-script/core/internal/store"
)

func TestDuplicate(t *testing.T) {
	chars := svccharacter.NewWithStore(store.NewMemory())
	voices := svcvoice.New(chars, allowAllPresets{})
	if err := chars.Add(domain.Character{Name: "A"}); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := voices.Add(ctx, "A", domain.Voice{
		Name: "main",
		Params: []domain.SynthParam{{
			Version:   domain.Version2,
			Speed:     100,
			PhontData: []byte{1, 2, 3},
		}},
	}); err != nil {
		t.Fatal(err)
	}
	if err := voices.Duplicate(ctx, "A", "main", ""); err != nil {
		t.Fatal(err)
	}
	got, err := voices.Get("A", "main - copy")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Params) != 1 || string(got.Params[0].PhontData) != "\x01\x02\x03" {
		t.Fatalf("%+v", got)
	}
	got.Params[0].PhontData[0] = 9
	orig, _ := voices.Get("A", "main")
	if orig.Params[0].PhontData[0] != 1 {
		t.Fatal("phont data not deep-copied")
	}
}
