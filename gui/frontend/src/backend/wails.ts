import * as App from '../../bindings/yukumo-gui/app';
import type { Backend, SynthRequest } from './backend';
import type { Character, Config, Voice } from './types';
import { base64ToBytes, bytesToBase64 } from './wav';

export class WailsBackend implements Backend {
  readonly kind = 'wails' as const;

  listCharacters(): Promise<Character[]> {
    return App.ListCharacters();
  }

  getCharacter(name: string): Promise<Character> {
    return App.GetCharacter(name);
  }

  addCharacter(character: Character): Promise<void> {
    return App.AddCharacter(character);
  }

  updateCharacter(oldName: string, character: Character): Promise<void> {
    return App.UpdateCharacter(oldName, character);
  }

  deleteCharacter(name: string): Promise<void> {
    return App.DeleteCharacter(name);
  }

  listVoices(charName: string): Promise<Voice[]> {
    return App.ListVoices(charName);
  }

  getVoice(charName: string, voiceName: string): Promise<Voice> {
    return App.GetVoice(charName, voiceName);
  }

  addVoice(charName: string, voice: Voice): Promise<void> {
    return App.AddVoice(charName, voice);
  }

  updateVoice(
    charName: string,
    oldName: string,
    voice: Voice,
  ): Promise<void> {
    return App.UpdateVoice(charName, oldName, voice);
  }

  deleteVoice(charName: string, voiceName: string): Promise<void> {
    return App.DeleteVoice(charName, voiceName);
  }

  duplicateVoice(
    charName: string,
    srcName: string,
    newName: string,
  ): Promise<string> {
    return App.DuplicateVoice(charName, srcName, newName);
  }

  config(): Promise<Config> {
    return App.Config();
  }

  saveConfig(config: Config): Promise<void> {
    return App.SaveConfig(config);
  }

  presetNames(version: number): Promise<string[]> {
    return App.PresetNames(version);
  }

  hasPreset(version: number, preset: number): Promise<boolean> {
    return App.HasPreset(version, preset);
  }

  async synthesize(request: SynthRequest): Promise<Uint8Array> {
    const encoded = await App.Synthesize({
      voice: request.voice,
      text: request.text,
      version: request.version ?? 0,
      speedFactor: request.speedFactor ?? 0,
    });
    return base64ToBytes(encoded);
  }

  async playWav(wav: Uint8Array): Promise<void> {
    await App.PlayWAV(bytesToBase64(wav));
  }

  async stopPlayback(): Promise<void> {
    await App.StopPlayback();
  }
}
