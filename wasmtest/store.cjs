'use strict';

// MemoryStore mirrors port.Store over a Map. It is the reference host adapter
// for the Node smoke test; a browser host can replace it with IndexedDB.
class MemoryStore {
  constructor(name) {
    this.name = name;
    this.data = new Map();
  }

  read(key) {
    const v = this.data.get(key);
    return v ? new Uint8Array(v) : null;
  }

  write(key, bytes) {
    this.data.set(key, new Uint8Array(bytes));
  }

  exists(key) {
    return this.data.has(key);
  }

  list(prefix) {
    const p = prefix || '';
    const out = [];
    for (const key of this.data.keys()) {
      if (!p) {
        if (!key.includes('/')) out.push(key);
        continue;
      }
      if (key.startsWith(p + '/')) {
        const rest = key.slice(p.length + 1);
        if (!rest.includes('/')) out.push(rest);
      }
    }
    return out.sort();
  }

  delete(key) {
    if (!this.data.has(key)) throw new Error('store: not found');
    this.data.delete(key);
  }

  path(key) {
    return `mem://${this.name}/${key}`;
  }

  flush() {
    if (globalThis.__yukumo) {
      globalThis.__yukumo.flushes = (globalThis.__yukumo.flushes || 0) + 1;
    }
  }
}

module.exports = { MemoryStore };
