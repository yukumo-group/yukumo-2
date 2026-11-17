'use strict';

// makeAqEngine exposes the port.Engine surface over the 0-import
// third-party/aq_wasm.wasm module. All calls are synchronous.
function makeAqEngine(instance) {
  const ex = instance.exports;
  const encoder = new TextEncoder();
  const decoder = new TextDecoder();

  const memory = () => new Uint8Array(ex.memory.buffer);
  const view = () => new DataView(ex.memory.buffer);

  function malloc(n) {
    const ptr = ex.malloc(n);
    if (!ptr && n) throw new Error(`malloc(${n}) returned null`);
    return ptr;
  }

  function free(ptr) {
    if (ptr) ex.free(ptr);
  }

  function readCString(ptr) {
    if (!ptr) return '';
    const mem = memory();
    let end = ptr;
    while (end < mem.length && mem[end] !== 0) end++;
    return decoder.decode(mem.subarray(ptr, end));
  }

  function errorMessage(version, code) {
    switch (code) {
      case -1: return 'AQ_ERR_ARG';
      case -2: return 'AQ_ERR_VERSION';
      case -3: return 'AQ_ERR_OPEN';
      default: {
        const ptr = ex.aq_error_message(version, code);
        return readCString(ptr) || 'unknown error';
      }
    }
  }

  function synthesize(text, param, phont) {
    const textBytes = encoder.encode(text);
    const textPtr = malloc(textBytes.length);
    memory().set(textBytes, textPtr);

    let phontPtr = 0;
    let phontLen = 0;
    if (phont && phont.length) {
      phontPtr = malloc(phont.length);
      memory().set(phont, phontPtr);
      phontLen = phont.length;
    }

    const paramPtr = malloc(40);
    const pv = new DataView(ex.memory.buffer, paramPtr, 40);
    pv.setInt32(0, param.version, true);
    pv.setInt32(4, param.preset, true);
    pv.setInt32(8, param.speed, true);
    pv.setUint32(12, phontPtr >>> 0, true);
    pv.setInt32(16, param.base, true);
    pv.setInt32(20, param.volume, true);
    pv.setInt32(24, param.pitch, true);
    pv.setInt32(28, param.accent, true);
    pv.setInt32(32, param.lmd, true);
    pv.setInt32(36, param.fsc, true);

    const outPtrSlot = malloc(4);
    const outLenSlot = malloc(4);
    const slots = view();
    slots.setUint32(outPtrSlot, 0, true);
    slots.setUint32(outLenSlot, 0, true);

    const code = ex.aq_synthesize(
      textPtr, textBytes.length, paramPtr, phontLen, outPtrSlot, outLenSlot,
    );
    if (code !== 0) {
      const msg = errorMessage(param.version, code);
      free(textPtr);
      free(phontPtr);
      free(paramPtr);
      free(outPtrSlot);
      free(outLenSlot);
      throw new Error(`aq_synthesize: ${msg} (code ${code})`);
    }

    const out = view();
    const wavPtr = out.getUint32(outPtrSlot, true);
    const wavLen = out.getUint32(outLenSlot, true);
    const wav = memory().slice(wavPtr, wavPtr + wavLen);
    ex.aq_free(wavPtr);

    free(textPtr);
    free(phontPtr);
    free(paramPtr);
    free(outPtrSlot);
    free(outLenSlot);
    return wav;
  }

  function hasPreset(version, preset) {
    return ex.aq_has_preset(version, preset) !== 0;
  }

  function presetNames(version) {
    const countPtr = malloc(4);
    view().setUint32(countPtr, 0, true);
    const tablePtr = ex.aq_preset_names(version, countPtr);
    if (!tablePtr) {
      free(countPtr);
      throw new Error(`unknown version ${version}`);
    }
    const count = view().getUint32(countPtr, true);
    const names = [];
    for (let i = 0; i < count; i++) {
      const ptr = view().getUint32(tablePtr + i * 4, true);
      names.push(readCString(ptr));
    }
    free(countPtr);
    return names;
  }

  return { synthesize, hasPreset, presetNames };
}

module.exports = { makeAqEngine };
