// check-wasm-deps inspects the GOOS=js dependency closure of wasmtest.
// It fails if runtime/wazero/cli/filesystem adapters enter the build, or if a
// project package imports os/path/filepath directly.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type pkg struct {
	ImportPath string
	Imports    []string
}

var forbidden = []string{
	"yukumo-script/runtime",
	"github.com/tetratelabs/wazero",
	"yukumo-script/cli",
	"yukumo-script/core/api/cliopen",
	"yukumo-script/core/internal/store/fs",
}

func main() {
	cmd := exec.Command("go", "list", "-deps", "-json", "./...")
	cmd.Dir = "wasmtest"
	cmd.Env = append(os.Environ(), "GOWORK=off", "GOOS=js", "GOARCH=wasm")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "go list: %v\n%s", err, stderr.String())
		os.Exit(1)
	}

	dec := json.NewDecoder(bytes.NewReader(out))
	failed := false
	for {
		var p pkg
		if err := dec.Decode(&p); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "decode go list: %v\n", err)
			os.Exit(1)
		}
		for _, bad := range forbidden {
			if p.ImportPath == bad || strings.HasPrefix(p.ImportPath, bad+"/") {
				fmt.Fprintf(os.Stderr, "wasm dependency %q is forbidden\n", p.ImportPath)
				failed = true
			}
		}
		if !strings.HasPrefix(p.ImportPath, "yukumo-script/") {
			continue
		}
		for _, imp := range p.Imports {
			if imp == "os" || imp == "path/filepath" {
				fmt.Fprintf(os.Stderr, "%s directly imports %s\n", p.ImportPath, imp)
				failed = true
			}
		}
	}
	if failed {
		os.Exit(1)
	}
	fmt.Println("wasm dependency closure: ok")
}
