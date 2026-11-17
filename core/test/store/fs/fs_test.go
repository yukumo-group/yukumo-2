package fs_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"yukumo-script/core/internal/store/fs"
	"yukumo-script/core/port"
)

func TestStoreAndOutput(t *testing.T) {
	st := fs.New(t.TempDir())
	if err := st.Write("chars/a.yaml", []byte("hello")); err != nil {
		t.Fatal(err)
	}
	got, err := st.Read("chars/a.yaml")
	if err != nil || string(got) != "hello" {
		t.Fatalf("read = %q, %v", got, err)
	}
	ok, err := st.Exists("chars/a.yaml")
	if err != nil || !ok {
		t.Fatalf("exists = %v, %v", ok, err)
	}
	list, err := st.List("chars")
	if err != nil || len(list) != 1 || list[0] != "a.yaml" {
		t.Fatalf("list = %v, %v", list, err)
	}
	if p := st.Path("chars/a.yaml"); !filepath.IsAbs(p) {
		t.Fatalf("path = %q", p)
	}
	if err := st.Delete("missing.yaml"); !errors.Is(err, port.ErrNotFound) {
		t.Fatalf("delete missing = %v", err)
	}
	if err := st.Delete("chars/a.yaml"); err != nil {
		t.Fatal(err)
	}
	if _, err := st.Read("chars/a.yaml"); !errors.Is(err, port.ErrNotFound) {
		t.Fatalf("read deleted = %v", err)
	}
	if err := st.Write("../escape", []byte("x")); err == nil {
		t.Fatal("expected key escape error")
	}

	outPath := filepath.Join(t.TempDir(), "nested", "out.wav")
	var out fs.Output
	if err := out.Write(outPath, []byte("wav")); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(outPath)
	if err != nil || string(data) != "wav" {
		t.Fatalf("output read = %q, %v", data, err)
	}
}
