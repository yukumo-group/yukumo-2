package aq

import (
	"context"
	"fmt"

	"github.com/tetratelabs/wazero/api"
)

// Synthesize runs aq_synthesize and returns a WAV byte slice.
// The Runtime must not be used concurrently from another goroutine.
func (r *Runtime) Synthesize(ctx context.Context, text string, p Param) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.mod == nil {
		return nil, fmt.Errorf("aq runtime closed")
	}

	malloc := r.mod.ExportedFunction("malloc")
	free := r.mod.ExportedFunction("free")
	synth := r.mod.ExportedFunction("aq_synthesize")
	aqFree := r.mod.ExportedFunction("aq_free")
	errMsg := r.mod.ExportedFunction("aq_error_message")
	if malloc == nil || free == nil || synth == nil || aqFree == nil {
		return nil, fmt.Errorf("aq synthesize exports missing")
	}

	mem := r.mod.Memory()
	textBytes := []byte(text)

	textPtr, err := mallocBytes(ctx, malloc, mem, textBytes)
	if err != nil {
		return nil, err
	}
	defer free.Call(ctx, uint64(textPtr)) //nolint:errcheck

	var phontPtr uint32
	var phontLen uint32
	if len(p.Phont) > 0 {
		phontPtr, err = mallocBytes(ctx, malloc, mem, p.Phont)
		if err != nil {
			return nil, err
		}
		defer free.Call(ctx, uint64(phontPtr)) //nolint:errcheck
		phontLen = uint32(len(p.Phont))
	}

	paramBytes := encodeParam(p, phontPtr)
	paramPtr, err := mallocBytes(ctx, malloc, mem, paramBytes)
	if err != nil {
		return nil, err
	}
	defer free.Call(ctx, uint64(paramPtr)) //nolint:errcheck

	outPtrSlot, err := mallocZero(ctx, malloc, mem, 4)
	if err != nil {
		return nil, err
	}
	defer free.Call(ctx, uint64(outPtrSlot)) //nolint:errcheck

	outLenSlot, err := mallocZero(ctx, malloc, mem, 4)
	if err != nil {
		return nil, err
	}
	defer free.Call(ctx, uint64(outLenSlot)) //nolint:errcheck

	results, err := synth.Call(ctx,
		uint64(textPtr),
		uint64(len(textBytes)),
		uint64(paramPtr),
		uint64(phontLen),
		uint64(outPtrSlot),
		uint64(outLenSlot),
	)
	if err != nil {
		return nil, err
	}
	code := api.DecodeI32(results[0])
	if code != 0 {
		msg := synthesizeErrorMessage(ctx, errMsg, mem, p.Version, code)
		return nil, fmt.Errorf("aq_synthesize failed: %s (code %d)", msg, code)
	}

	wavPtr, ok := mem.ReadUint32Le(outPtrSlot)
	if !ok {
		return nil, fmt.Errorf("read wav out_ptr")
	}
	wavLen, ok := mem.ReadUint32Le(outLenSlot)
	if !ok {
		return nil, fmt.Errorf("read wav out_len")
	}
	if wavPtr == 0 || wavLen == 0 {
		return nil, fmt.Errorf("aq_synthesize returned empty wav")
	}
	defer aqFree.Call(ctx, uint64(wavPtr)) //nolint:errcheck

	wav, ok := mem.Read(wavPtr, wavLen)
	if !ok {
		return nil, fmt.Errorf("read wav bytes")
	}
	out := make([]byte, len(wav))
	copy(out, wav)
	return out, nil
}

func mallocBytes(ctx context.Context, malloc api.Function, mem api.Memory, data []byte) (uint32, error) {
	if len(data) == 0 {
		return 0, nil
	}
	results, err := malloc.Call(ctx, uint64(len(data)))
	if err != nil {
		return 0, err
	}
	ptr := uint32(results[0])
	if ptr == 0 {
		return 0, fmt.Errorf("malloc %d returned null", len(data))
	}
	if !mem.Write(ptr, data) {
		return 0, fmt.Errorf("write malloc buffer")
	}
	return ptr, nil
}

func mallocZero(ctx context.Context, malloc api.Function, mem api.Memory, n uint32) (uint32, error) {
	results, err := malloc.Call(ctx, uint64(n))
	if err != nil {
		return 0, err
	}
	ptr := uint32(results[0])
	if ptr == 0 {
		return 0, fmt.Errorf("malloc %d returned null", n)
	}
	zeros := make([]byte, n)
	if !mem.Write(ptr, zeros) {
		return 0, fmt.Errorf("zero malloc buffer")
	}
	return ptr, nil
}

// Host-side codes from aq_wasm.h (returned before an engine runs).
const (
	aqErrArg     int32 = -1 // AQ_ERR_ARG
	aqErrVersion int32 = -2 // AQ_ERR_VERSION
	aqErrOpen    int32 = -3 // AQ_ERR_OPEN
)

func synthesizeErrorMessage(ctx context.Context, errMsg api.Function, mem api.Memory, version, code int32) string {
	switch code {
	case aqErrArg:
		return "AQ_ERR_ARG"
	case aqErrVersion:
		return "AQ_ERR_VERSION"
	case aqErrOpen:
		return "AQ_ERR_OPEN"
	}
	if s := lookupEngineError(ctx, errMsg, mem, version, code); s != "" {
		return s
	}
	return "unknown error"
}

func lookupEngineError(ctx context.Context, errMsg api.Function, mem api.Memory, version, code int32) string {
	if errMsg == nil {
		return ""
	}
	results, err := errMsg.Call(ctx, api.EncodeI32(version), api.EncodeI32(code))
	if err != nil || len(results) == 0 {
		return ""
	}
	ptr := uint32(results[0])
	if ptr == 0 {
		return ""
	}
	s, err := readCString(mem, ptr)
	if err != nil {
		return ""
	}
	return s
}
