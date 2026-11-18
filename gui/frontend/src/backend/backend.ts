import type { Character, Config, Voice } from './types';

export interface SynthRequest {
  voice: Voice;
  text: string;
  version?: number;
  speedFactor?: number;
}

export interface Backend {
  readonly kind: 'wails' | 'wasm';
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
  duplicateVoice(
    charName: string,
    srcName: string,
    newName: string,
  ): Promise<string>;
  config(): Promise<Config>;
  saveConfig(config: Config): Promise<void>;
  presetNames(version: number): Promise<string[]>;
  hasPreset(version: number, preset: number): Promise<boolean>;
  synthesize(request: SynthRequest): Promise<Uint8Array>;
  playWav(wav: Uint8Array): Promise<void>;
  stopPlayback(): Promise<void>;
}
