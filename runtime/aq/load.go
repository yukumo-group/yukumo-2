package aq

import (
	"fmt"
	"os"
	"path/filepath"
)

func loadWASM() ([]byte, error) {
	for _, p := range wasmCandidates() {
		data, err := os.ReadFile(p)
		if err == nil {
			return data, nil
		}
	}
	return nil, fmt.Errorf("aq_wasm.wasm not found (set YUKUMO_AQ_WASM or place third-party/aq_wasm.wasm)")
}

func wasmCandidates() []string {
	var out []string
	if p := os.Getenv("YUKUMO_AQ_WASM"); p != "" {
		out = append(out, p)
	}
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		out = append(out,
			filepath.Join(dir, "aq_wasm.wasm"),
			filepath.Join(dir, "third-party", "aq_wasm.wasm"),
		)
	}
	wd, err := os.Getwd()
	if err != nil {
		return out
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		out = append(out, filepath.Join(dir, "third-party", "aq_wasm.wasm"))
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}
	return out
}
