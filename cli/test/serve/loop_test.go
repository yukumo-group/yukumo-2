package serve_test

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"yukumo-script/cli/serve"
	"yukumo-script/core/api/cliopen"
	"yukumo-script/core/api/rpc"
)

func TestLoopCharacterList(t *testing.T) {
	dir := t.TempDir()
	app, err := cliopen.Open(cliopen.Options{
		ConfigDir:     dir,
		CharacterPath: filepath.Join(dir, "characters"),
		OutputPath:    filepath.Join(dir, "output"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = app.Close(context.Background()) })

	in := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"character.list","params":{}}` + "\n")
	var out bytes.Buffer
	var errBuf bytes.Buffer
	if err := serve.Loop(context.Background(), rpc.Bind(app), in, &out, &errBuf); err != nil {
		t.Fatal(err)
	}
	line := strings.TrimSpace(out.String())
	var res map[string]any
	if err := json.Unmarshal([]byte(line), &res); err != nil {
		t.Fatalf("unmarshal %q: %v", line, err)
	}
	if res["error"] != nil {
		t.Fatalf("error: %v", res["error"])
	}
	if errBuf.Len() != 0 {
		t.Fatalf("stderr: %s", errBuf.String())
	}
}
