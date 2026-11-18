import type { Backend, SynthRequest } from './backend';
import type { Character, Config, Voice } from './types';
import { base64ToBytes, createPlayback } from './wav';

export interface WasmBackendHost {
  listCharacters(): Promise<Character[]>;
  getCharacter(name: string): Promise<Character>;
  addCharacter(character: Character): Promise<void>;
  updateCharacter(oldName: string, character: Character): Promise<void>;
  deleteCharacter(name: string): Promise<void>;
  listVoices(charName: string): Promise<Voice[]>;
  getVoice(charName: string, voiceName: string): Promise<Voice>;
  addVoice(charName: string, voice: Voice): Promise<void>;
  updateVoice(charName: string, oldName: string, voice: Voice): Promise<void>;
  deleteVoice(charName: string, voiceName: string): Promise<void>;
  duplicateVoice(charName: string, srcName: string, newName: string): Promise<string>;
  config(): Promise<Config>;
  saveConfig(config: Config): Promise<void>;
  presetNames(version: number): Promise<string[]>;
  hasPreset(version: number, preset: number): Promise<boolean>;
  synthesize(request: {
    voice: Voice;
    text: string;
    version: number;
    speedFactor: number;
  }): Promise<string>;
}

declare global {
  interface Window {
    __yukumoBackend?: WasmBackendHost;
  }
}

function wasmHost(): WasmBackendHost {
  const host = globalThis.window?.__yukumoBackend;
  if (!host) {
    throw new Error('WASM backend is not available in this host');
  }
  return host;
}

export class WasmBackend implements Backend {
  readonly kind = 'wasm' as const;

  private stopAudio: (() => void) | null = null;

  listCharacters(): Promise<Character[]> {
    return wasmHost().listCharacters();
  }

  getCharacter(name: string): Promise<Character> {
    return wasmHost().getCharacter(name);
  }

  addCharacter(character: Character): Promise<void> {
    return wasmHost().addCharacter(character);
  }

  updateCharacter(oldName: string, character: Character): Promise<void> {
    return wasmHost().updateCharacter(oldName, character);
  }

  deleteCharacter(name: string): Promise<void> {
    return wasmHost().deleteCharacter(name);
  }

  listVoices(charName: string): Promise<Voice[]> {
    return wasmHost().listVoices(charName);
  }

  getVoice(charName: string, voiceName: string): Promise<Voice> {
    return wasmHost().getVoice(charName, voiceName);
  }

  addVoice(charName: string, voice: Voice): Promise<void> {
    return wasmHost().addVoice(charName, voice);
  }

  updateVoice(
    charName: string,
    oldName: string,
    voice: Voice,
  ): Promise<void> {
    return wasmHost().updateVoice(charName, oldName, voice);
  }

  deleteVoice(charName: string, voiceName: string): Promise<void> {
    return wasmHost().deleteVoice(charName, voiceName);
  }

  duplicateVoice(
    charName: string,
    srcName: string,
    newName: string,
  ): Promise<string> {
    return wasmHost().duplicateVoice(charName, srcName, newName);
  }

  config(): Promise<Config> {
    return wasmHost().config();
  }

  saveConfig(config: Config): Promise<void> {
    return wasmHost().saveConfig(config);
  }

  presetNames(version: number): Promise<string[]> {
    return wasmHost().presetNames(version);
  }

  hasPreset(version: number, preset: number): Promise<boolean> {
    return wasmHost().hasPreset(version, preset);
  }

  async synthesize(request: SynthRequest): Promise<Uint8Array> {
    const encoded = await wasmHost().synthesize({
      voice: request.voice,
      text: request.text,
      version: request.version ?? 0,
      speedFactor: request.speedFactor ?? 0,
    });
    return base64ToBytes(encoded);
  }

  async playWav(wav: Uint8Array): Promise<void> {
    this.stopPlaybackNow();
    const playback = createPlayback(wav, () => {
      this.stopAudio = null;
    });
    this.stopAudio = playback.stop;
    await playback.audio.play();
  }

  async stopPlayback(): Promise<void> {
    this.stopPlaybackNow();
  }

  private stopPlaybackNow(): void {
    this.stopAudio?.();
    this.stopAudio = null;
  }
}
