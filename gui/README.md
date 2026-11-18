Moved from [yukumo-gui](https://github.com/yukumo-group/yukumo-gui).

# yukumo-gui

Desktop GUI for [yukumo-script](https://github.com/yukumo-group/yukumo-script), built with [Wails v3](https://v3.wails.io/) and Vue 3.

The GUI lives inside the `yukumo-script` repository and uses the local Go modules through `replace` directives in `go.mod`: `yukumo-script/core/api`, `core`, `core/port`, and `runtime`. There are no git submodules.

The frontend talks to the backend through `frontend/src/backend/`:

- `WailsBackend` calls the generated Wails bindings for the native app.
- `WasmBackend` is the seam for a future Go `js/wasm` host (`window.__yukumoBackend`).

Both adapters implement the same `Backend` interface, so UI code does not depend on Wails.

## Live Development

```bash
wails3 dev
```

This starts the app with hot-reload for frontend changes. Vite listens on port `9245` by default (`WAILS_VITE_PORT`).
On systems that only provide GTK3/WebKitGTK 4.1, build with the `gtk3` tag:

```bash
wails3 build -tags gtk3
```

## Building

```bash
wails3 build
```

Production binaries are written under `bin/`. Frontend dependencies are installed with **pnpm** from `frontend/`.

The native engine loads `third-party/aq_wasm.wasm` from the executable directory, its parent, or the repository root. For a packaged build, set `YUKUMO_AQ_WASM` or copy `aq_wasm.wasm` next to the binary.

## Linux / Wayland troubleshooting

Preview audio uses a native OS player (`paplay`, `aplay`, `ffplay`, or `mpv`), so WebKitGTK's GStreamer `autoaudiosink` is not required.

If the window is sluggish, stalls, or stops repainting on Wayland or a hybrid-GPU setup, try the WebKitGTK compatibility switches:

```bash
WEBKIT_DISABLE_DMABUF_RENDERER=1 wails3 dev
GDK_BACKEND=x11 wails3 dev
```

GTK4 builds can additionally try `GSK_RENDERER=cairo`.

## Bindings

After changing the Wails service surface in `app.go`, regenerate the TypeScript bindings:

```bash
wails3 generate bindings -ts -clean -d frontend/bindings -f '-tags=gtk3'
```
