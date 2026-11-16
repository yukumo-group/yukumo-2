package aq

import (
	"context"
	"fmt"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type Runtime struct {
	mu   sync.Mutex
	rt   wazero.Runtime
	mod  api.Module
	free func(context.Context) error
}

func New(ctx context.Context) (*Runtime, error) {
	wasm, err := loadWASM()
	if err != nil {
		return nil, err
	}
	rt := wazero.NewRuntime(ctx)
	mod, err := rt.Instantiate(ctx, wasm)
	if err != nil {
		_ = rt.Close(ctx)
		return nil, fmt.Errorf("instantiate aq_wasm: %w", err)
	}
	return &Runtime{
		rt:  rt,
		mod: mod,
		free: func(c context.Context) error {
			cerr := mod.Close(c)
			rerr := rt.Close(c)
			if cerr != nil {
				return cerr
			}
			return rerr
		},
	}, nil
}

func (r *Runtime) Close(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.free == nil {
		return nil
	}
	err := r.free(ctx)
	r.free = nil
	r.mod = nil
	r.rt = nil
	return err
}

func (r *Runtime) HasPreset(ctx context.Context, version, preset int32) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	fn := r.mod.ExportedFunction("aq_has_preset")
	if fn == nil {
		return false, fmt.Errorf("aq_has_preset not exported")
	}
	results, err := fn.Call(ctx, api.EncodeI32(version), api.EncodeI32(preset))
	if err != nil {
		return false, err
	}
	return api.DecodeI32(results[0]) != 0, nil
}

func (r *Runtime) PresetNames(ctx context.Context, version int32) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	malloc := r.mod.ExportedFunction("malloc")
	free := r.mod.ExportedFunction("free")
	namesFn := r.mod.ExportedFunction("aq_preset_names")
	if malloc == nil || free == nil || namesFn == nil {
		return nil, fmt.Errorf("aq preset exports missing")
	}

	countPtrResults, err := malloc.Call(ctx, 4)
	if err != nil {
		return nil, err
	}
	countPtr := uint32(countPtrResults[0])
	defer free.Call(ctx, uint64(countPtr)) //nolint:errcheck

	mem := r.mod.Memory()
	if !mem.WriteUint32Le(countPtr, 0) {
		return nil, fmt.Errorf("write out_count")
	}

	results, err := namesFn.Call(ctx, api.EncodeI32(version), uint64(countPtr))
	if err != nil {
		return nil, err
	}
	tablePtr := uint32(results[0])
	if tablePtr == 0 {
		return nil, fmt.Errorf("unknown version %d", version)
	}
	count, ok := mem.ReadUint32Le(countPtr)
	if !ok {
		return nil, fmt.Errorf("read out_count")
	}

	out := make([]string, 0, count)
	for i := uint32(0); i < count; i++ {
		strPtr, ok := mem.ReadUint32Le(tablePtr + i*4)
		if !ok {
			return nil, fmt.Errorf("read name ptr %d", i)
		}
		s, err := readCString(mem, strPtr)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

func readCString(mem api.Memory, ptr uint32) (string, error) {
	if ptr == 0 {
		return "", nil
	}
	var b []byte
	for {
		c, ok := mem.ReadByte(ptr)
		if !ok {
			return "", fmt.Errorf("read cstring at %d", ptr)
		}
		if c == 0 {
			return string(b), nil
		}
		b = append(b, c)
		ptr++
	}
}
