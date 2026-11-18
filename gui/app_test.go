package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
	coreapi "yukumo-script/core/api"
	"yukumo-script/core/api/cliopen"
)

func TestCharacterCRUDAndInlineSynth(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	app := &App{openOptions: cliopen.Options{
		ConfigDir:     filepath.Join(dir, "config"),
		CharacterPath: filepath.Join(dir, "characters"),
		OutputPath:    filepath.Join(dir, "output"),
	}}
	if err := app.ServiceStartup(ctx, application.ServiceOptions{}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = app.ServiceShutdown() })

	character := coreapi.Character{
		Name:        "gui-test",
		Description: "created from gui",
		Voices: []coreapi.Voice{{
			Name:   "default",
			Params: []coreapi.SynthParam{{Version: coreapi.Version1, Preset: 0, Speed: 100}},
		}},
	}
	if err := app.AddCharacter(character); err != nil {
		t.Fatal(err)
	}
	list, err := app.ListCharacters()
	if err != nil || len(list) != 1 || list[0].Name != "gui-test" {
		t.Fatalf("list = %+v, %v", list, err)
	}

	character.Description = "updated"
	if err := app.UpdateCharacter("gui-test", character); err != nil {
		t.Fatal(err)
	}
	got, err := app.GetCharacter("gui-test")
	if err != nil || got.Description != "updated" {
		t.Fatalf("get = %+v, %v", got, err)
	}

	names, err := app.PresetNames(ctx, coreapi.Version1)
	if err != nil || len(names) == 0 {
		t.Fatalf("preset names = %v, %v", names, err)
	}

	wav, err := app.Synthesize(ctx, SynthInput{
		Voice: coreapi.Voice{
			Name:   "inline",
			Params: []coreapi.SynthParam{{Version: coreapi.Version1, Preset: 0, Speed: 100}},
		},
		Text:    "こんにちは",
		Version: coreapi.Version1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(wav) < 44 || string(wav[:4]) != "RIFF" {
		t.Fatalf("bad wav: len=%d", len(wav))
	}

	entries, err := os.ReadDir(filepath.Join(dir, "output"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("synthesize must not write artifacts, got %v", entries)
	}

	if err := app.DeleteCharacter("gui-test"); err != nil {
		t.Fatal(err)
	}
	if _, err := app.GetCharacter("gui-test"); err == nil {
		t.Fatal("character still exists after delete")
	}
}
