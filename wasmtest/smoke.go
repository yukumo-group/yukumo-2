//go:build js

package main

import (
	"bytes"
	"context"
	"fmt"
	"syscall/js"

	"yukumo-script/core/api"
)

func runSmoke() error {
	ctx := context.Background()
	configStore := newJSStore(storeRef("config"))
	charsStore := newJSStore(storeRef("characters"))
	artifacts := newJSArtifactStore(storeRef("artifacts"))

	if err := fakeScenario(ctx, configStore, charsStore, artifacts); err != nil {
		return fmt.Errorf("fake engine: %w", err)
	}
	if err := realScenario(ctx, configStore, charsStore, artifacts); err != nil {
		return fmt.Errorf("real engine: %w", err)
	}
	if js.Global().Get("__yukumo").Get("flushes").Int() == 0 {
		return fmt.Errorf("stores were not flushed")
	}
	return nil
}

func fakeScenario(ctx context.Context, cfgStore, charsStore JSStore, artifacts JSArtifactStore) error {
	app, err := api.New(api.Deps{
		ConfigStore:    cfgStore,
		CharacterStore: charsStore,
		Artifacts:      artifacts,
		Engine:         fakeEngine{},
		Defaults:       api.DefaultConfig(),
	})
	if err != nil {
		return err
	}
	if err := app.SaveConfig(api.DefaultConfig()); err != nil {
		return err
	}
	if _, err := app.Config(); err != nil {
		return err
	}
	if err := app.AddCharacter(api.Character{Name: "smoke", Description: "fake"}); err != nil {
		return err
	}
	if err := app.AddVoice(ctx, "smoke", api.Voice{
		Name:   "default",
		Params: []api.SynthParam{{Version: api.Version1, Preset: 0, Speed: 100}},
	}); err != nil {
		return err
	}
	voices, err := app.ListVoices("smoke")
	if err != nil {
		return err
	}
	if len(voices) != 1 || voices[0].Name != "default" {
		return fmt.Errorf("unexpected voices: %+v", voices)
	}
	if _, err := app.DuplicateVoice(ctx, "smoke", "default", "copy"); err != nil {
		return err
	}
	if err := app.DeleteVoice("smoke", "copy"); err != nil {
		return err
	}
	return app.Close(ctx)
}

func realScenario(ctx context.Context, cfgStore, charsStore JSStore, artifacts JSArtifactStore) error {
	app, err := api.New(api.Deps{
		ConfigStore:    cfgStore,
		CharacterStore: charsStore,
		Artifacts:      artifacts,
		Engine:         newJSEngine(js.Global().Get("__yukumo").Get("engine")),
		Defaults:       api.DefaultConfig(),
	})
	if err != nil {
		return err
	}
	defer func() { _ = app.Close(context.Background()) }()

	if err := app.AddCharacter(api.Character{Name: "aqsmoke"}); err != nil {
		return err
	}
	if err := app.AddVoice(ctx, "aqsmoke", api.Voice{
		Name:   "default",
		Params: []api.SynthParam{{Version: api.Version1, Preset: 0, Speed: 100}},
	}); err != nil {
		return err
	}
	first, err := app.Synthesize(ctx, api.SynthRequest{
		VoiceRef: "aqsmoke:default",
		Text:     "こんにちは",
	})
	if err != nil {
		return err
	}
	if !first.Wrote || first.Path == "" {
		return fmt.Errorf("first synth did not write: %+v", first)
	}
	if len(first.WAV) < 44 || string(first.WAV[:4]) != "RIFF" {
		return fmt.Errorf("bad wav: len=%d", len(first.WAV))
	}

	second, err := app.Synthesize(ctx, api.SynthRequest{
		VoiceRef: "aqsmoke:default",
		Text:     "こんにちは",
	})
	if err != nil {
		return err
	}
	if !second.FromCache {
		return fmt.Errorf("second synth was not cached: %+v", second)
	}
	if !bytes.Equal(first.WAV, second.WAV) {
		return fmt.Errorf("cached wav differs")
	}
	if err := app.CleanIndex(); err != nil {
		return err
	}
	return nil
}

func storeRef(name string) js.Value {
	return js.Global().Get("__yukumo").Get("stores").Get(name)
}
