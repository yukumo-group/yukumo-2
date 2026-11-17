package api

import (
	"reflect"
	"testing"

	"yukumo-script/core/internal/domain"
)

func TestCharacterRoundTripDeepCopies(t *testing.T) {
	in := domain.Character{
		Name:        "remilia",
		Description: "vampire",
		Voices: []domain.Voice{{
			Name: "default",
			Params: []domain.SynthParam{{
				Version:   domain.Version2,
				Preset:    1,
				Speed:     100,
				PhontData: []byte{1, 2, 3},
			}},
		}},
	}
	out := fromCharacter(toCharacter(in))
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("round trip mismatch:\nin=%+v\nout=%+v", in, out)
	}
	out.Voices[0].Params[0].PhontData[0] = 9
	if in.Voices[0].Params[0].PhontData[0] != 1 {
		t.Fatal("DTO conversion should deep-copy phont data")
	}
}

func TestConfigEffectiveBooleans(t *testing.T) {
	cfg := fromConfig(Config{
		Paths: PathConfig{CharacterPath: "chars", OutputPath: "out"},
		CLI:   CLIConfig{DefaultInteractive: true, WithAccessibility: false},
		Synth: SynthConfig{OutputFormat: "{character}", AlwaysWrite: false, ReuseCache: true},
	})
	if !cfg.CLI.Interactive() || cfg.CLI.Accessibility() {
		t.Fatal("CLI booleans not mapped")
	}
	if cfg.Synth.WriteAlways() || !cfg.Synth.CacheReuse() {
		t.Fatal("synth booleans not mapped")
	}
	back := toConfig(cfg)
	if back.Paths.CharacterPath != "chars" || back.Paths.OutputPath != "out" {
		t.Fatalf("paths not mapped: %+v", back.Paths)
	}
	if back.CLI.DefaultInteractive != true || back.Synth.AlwaysWrite != false {
		t.Fatalf("effective booleans not mapped: %+v", back)
	}
}

func TestDefaultSynthParamAndBaseNames(t *testing.T) {
	p := DefaultSynthParam(Version10)
	if p.Version != Version10 || p.Volume != 100 {
		t.Fatalf("unexpected default: %+v", p)
	}
	p.Speed = 999
	p.Normalize()
	if p.Speed != 300 {
		t.Fatalf("normalize speed = %d", p.Speed)
	}
	names := BaseNames()
	if len(names) == 0 || names[0] != "F1E" {
		t.Fatalf("unexpected base names: %v", names)
	}
	names[0] = "mutated"
	if BaseNames()[0] != "F1E" {
		t.Fatal("BaseNames must return a copy")
	}
}
