//go:build js

package main

import (
	"context"
	"fmt"
	"syscall/js"

	"yukumo-script/core/port"
)

// JSEngine adapts the host aq_wasm engine to port.Engine.
type JSEngine struct {
	ref js.Value
}

func newJSEngine(ref js.Value) JSEngine { return JSEngine{ref: ref} }

func (e JSEngine) Synthesize(_ context.Context, text string, p port.EngineParam) (out []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("js engine: %v", r)
		}
	}()
	param := js.Global().Get("Object").New()
	param.Set("version", p.Version)
	param.Set("preset", p.Preset)
	param.Set("speed", p.Speed)
	param.Set("base", p.Base)
	param.Set("volume", p.Volume)
	param.Set("pitch", p.Pitch)
	param.Set("accent", p.Accent)
	param.Set("lmd", p.Lmd)
	param.Set("fsc", p.Fsc)

	var phont js.Value
	if len(p.Phont) > 0 {
		phont = js.Global().Get("Uint8Array").New(len(p.Phont))
		js.CopyBytesToJS(phont, p.Phont)
	} else {
		phont = js.Null()
	}

	res := e.ref.Call("synthesize", text, param, phont)
	return copyBytesFromJS(res), nil
}

func (e JSEngine) HasPreset(_ context.Context, version, preset int32) (bool, error) {
	return e.ref.Call("hasPreset", int(version), int(preset)).Bool(), nil
}

func (e JSEngine) PresetNames(_ context.Context, version int32) ([]string, error) {
	v := e.ref.Call("presetNames", int(version))
	n := v.Length()
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = v.Index(i).String()
	}
	return out, nil
}
