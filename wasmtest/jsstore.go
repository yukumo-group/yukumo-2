//go:build js

package main

import (
	"context"
	"syscall/js"

	"yukumo-script/core/port"
)

// JSStore adapts the host MemoryStore (Map-backed) to port.Store.
type JSStore struct {
	ref js.Value
}

func newJSStore(ref js.Value) JSStore { return JSStore{ref: ref} }

func (s JSStore) Read(key string) ([]byte, error) {
	v := s.ref.Call("read", key)
	if v.IsNull() || v.IsUndefined() {
		return nil, port.ErrNotFound
	}
	return copyBytesFromJS(v), nil
}

func (s JSStore) Write(key string, data []byte) error {
	arr := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(arr, data)
	s.ref.Call("write", key, arr)
	return nil
}

func (s JSStore) Exists(key string) (bool, error) {
	return s.ref.Call("exists", key).Bool(), nil
}

func (s JSStore) List(prefix string) ([]string, error) {
	v := s.ref.Call("list", prefix)
	n := v.Length()
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = v.Index(i).String()
	}
	return out, nil
}

func (s JSStore) Delete(key string) error {
	ok := s.ref.Call("exists", key).Bool()
	if !ok {
		return port.ErrNotFound
	}
	s.ref.Call("delete", key)
	return nil
}

func (s JSStore) Flush(context.Context) error {
	s.ref.Call("flush")
	return nil
}

// JSArtifactStore adds the display-only Path method.
type JSArtifactStore struct {
	JSStore
}

func newJSArtifactStore(ref js.Value) JSArtifactStore {
	return JSArtifactStore{JSStore: newJSStore(ref)}
}

func (s JSArtifactStore) Path(key string) string {
	return s.ref.Call("path", key).String()
}

func copyBytesFromJS(v js.Value) []byte {
	n := v.Length()
	out := make([]byte, n)
	js.CopyBytesToGo(out, v)
	return out
}
