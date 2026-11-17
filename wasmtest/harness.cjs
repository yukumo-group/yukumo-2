'use strict';

const fs = require('node:fs');
const path = require('node:path');

// Match the host environment that GOROOT/lib/wasm/wasm_exec.js expects.
globalThis.require = require;
globalThis.fs = fs;
globalThis.path = path;
globalThis.TextEncoder ??= require('node:util').TextEncoder;
globalThis.TextDecoder ??= require('node:util').TextDecoder;
globalThis.performance ??= require('node:perf_hooks').performance;
globalThis.crypto ??= require('node:crypto').webcrypto;

require('./.build/wasm_exec.js');

const { MemoryStore } = require('./store.cjs');
const { makeAqEngine } = require('./aq_engine.cjs');

async function main() {
  const aqBytes = fs.readFileSync(path.join(__dirname, '..', 'third-party', 'aq_wasm.wasm'));
  const { instance: aqInstance } = await WebAssembly.instantiate(aqBytes);
  const engine = makeAqEngine(aqInstance);

  const config = new MemoryStore('config');
  const characters = new MemoryStore('characters');
  const artifacts = new MemoryStore('artifacts');

  globalThis.__yukumo = {
    engine,
    stores: { config, characters, artifacts },
    flushes: 0,
  };

  const go = new globalThis.Go();
  go.argv = [];
  go.env = { ...process.env };
  const program = fs.readFileSync(path.join(__dirname, '.build', 'smoke.wasm'));
  const { instance: goInstance } = await WebAssembly.instantiate(program, go.importObject);
  await go.run(goInstance);
}

main().catch((err) => {
  console.error(err);
  process.exitCode = 1;
});
