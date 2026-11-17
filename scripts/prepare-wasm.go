// prepare-wasm copies the Go version's wasm_exec.js into wasmtest/.build so
// the Node harness always matches the toolchain that built smoke.wasm.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	out, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "go env GOROOT: %v\n", err)
		os.Exit(1)
	}
	goroot := strings.TrimSpace(string(out))
	src := filepath.Join(goroot, "lib", "wasm", "wasm_exec.js")
	dst := filepath.Join("wasmtest", ".build", "wasm_exec.js")
	data, err := os.ReadFile(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", src, err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", filepath.Dir(dst), err)
		os.Exit(1)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", dst, err)
		os.Exit(1)
	}
}
