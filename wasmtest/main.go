//go:build js

package main

import "syscall/js"

func main() {
	if err := runSmoke(); err != nil {
		js.Global().Get("console").Call("error", "wasm smoke failed: "+err.Error())
		js.Global().Get("process").Set("exitCode", 1)
		return
	}
	js.Global().Get("console").Call("log", "wasm smoke: ok")
	js.Global().Get("process").Set("exitCode", 0)
}
