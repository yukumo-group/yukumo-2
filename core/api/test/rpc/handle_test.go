package rpc_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"yukumo-script/core/api"
	"yukumo-script/core/api/rpc"
	"yukumo-script/core/internal/store"
)

func testApp(t *testing.T) *api.App {
	t.Helper()
	app, err := api.New(api.Deps{
		ConfigStore:    store.NewMemory(),
		CharacterStore: store.NewMemory(),
		Artifacts:      store.NewMemory(),
		Defaults:       api.DefaultConfig(),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = app.Close(context.Background()) })
	return app
}

func TestHandleCharacterListAndAdd(t *testing.T) {
	reg := rpc.Bind(testApp(t))
	raw := reg.Handle(context.Background(), []byte(`{"jsonrpc":"2.0","id":1,"method":"character.list","params":{}}`))
	var res map[string]any
	if err := json.Unmarshal(raw, &res); err != nil {
		t.Fatal(err)
	}
	if res["error"] != nil {
		t.Fatalf("unexpected error: %v", res["error"])
	}
	result := res["result"].(map[string]any)
	chars := result["characters"].([]any)
	if len(chars) != 0 {
		t.Fatalf("want empty list, got %v", chars)
	}

	add := `{"jsonrpc":"2.0","id":2,"method":"character.add","params":{"character":{"name":"remilia","description":"x"}}}`
	raw = reg.Handle(context.Background(), []byte(add))
	if err := json.Unmarshal(raw, &res); err != nil {
		t.Fatal(err)
	}
	if res["error"] != nil {
		t.Fatalf("add error: %v", res["error"])
	}

	raw = reg.Handle(context.Background(), []byte(`{"jsonrpc":"2.0","id":3,"method":"character.list","params":{}}`))
	if err := json.Unmarshal(raw, &res); err != nil {
		t.Fatal(err)
	}
	result = res["result"].(map[string]any)
	chars = result["characters"].([]any)
	if len(chars) != 1 {
		t.Fatalf("want 1 character, got %v", chars)
	}
}

func TestHandleUnknownMethod(t *testing.T) {
	reg := rpc.Bind(testApp(t))
	raw := reg.Handle(context.Background(), []byte(`{"jsonrpc":"2.0","id":1,"method":"nope","params":{}}`))
	var res struct {
		Error *struct {
			Code int            `json:"code"`
			Data map[string]any `json:"data"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		t.Fatal(err)
	}
	if res.Error == nil || res.Error.Code != -32601 {
		t.Fatalf("want -32601, got %+v", res.Error)
	}
}

func TestHandleInvalidParams(t *testing.T) {
	reg := rpc.Bind(testApp(t))
	raw := reg.Handle(context.Background(), []byte(`{"jsonrpc":"2.0","id":1,"method":"character.get","params":[]}`))
	var res struct {
		Error *struct {
			Code int `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		t.Fatal(err)
	}
	if res.Error == nil || res.Error.Code != -32602 {
		t.Fatalf("want -32602, got %+v", res.Error)
	}
}

func TestHandleNotFoundKind(t *testing.T) {
	reg := rpc.Bind(testApp(t))
	raw := reg.Handle(context.Background(), []byte(`{"jsonrpc":"2.0","id":1,"method":"character.get","params":{"name":"missing"}}`))
	var res struct {
		Error *struct {
			Code int            `json:"code"`
			Data map[string]any `json:"data"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &res); err != nil {
		t.Fatal(err)
	}
	if res.Error == nil || res.Error.Code != -32000 {
		t.Fatalf("want -32000, got %+v", res.Error)
	}
	if res.Error.Data["kind"] != "not_found" {
		t.Fatalf("want kind not_found, got %v", res.Error.Data)
	}
}

func TestHandleParseError(t *testing.T) {
	reg := rpc.Bind(testApp(t))
	raw := reg.Handle(context.Background(), []byte(`{`))
	if !strings.Contains(string(raw), `"code":-32700`) {
		t.Fatalf("want parse error, got %s", raw)
	}
}

func TestBindMethodCount(t *testing.T) {
	reg := rpc.Bind(testApp(t))
	if n := len(reg.Methods()); n < 15 {
		t.Fatalf("expected at least 15 methods, got %d", n)
	}
}
