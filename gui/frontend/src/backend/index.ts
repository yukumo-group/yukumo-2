import type { Backend } from './backend';
import { WailsBackend } from './wails';
import { WasmBackend } from './wasm';

let backend: Backend | null = null;

function hasWailsRuntime(): boolean {
  const wails = (globalThis.window as Window & { _wails?: { invoke?: unknown } })
    ?._wails;
  return typeof wails?.invoke === 'function';
}

function createBackend(): Backend {
  return hasWailsRuntime() ? new WailsBackend() : new WasmBackend();
}

export function getBackend(): Backend {
  if (backend === null) {
    backend = createBackend();
  }
  return backend;
}

export function setBackend(next: Backend): void {
  backend = next;
}
