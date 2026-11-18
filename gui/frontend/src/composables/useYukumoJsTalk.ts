import { onBeforeUnmount, ref } from 'vue';
import type { AqKanji2Koe } from 'yukumo.js/lang/kanji2koe';
import {
  createPlayback,
  downloadWav,
  wavFilename,
} from '../yukumo-js/audio';
import {
  engineCacheKey,
  loadKanji2Koe,
  loadTalkEngine,
  synthesize,
  type TalkEngine,
} from '../yukumo-js/engines';
import type {
  AquesTalk1Config,
  AquesTalk10Config,
  AquesTalk2Config,
  AquesTalkVersion,
} from '../types/profile';

export type YukumoJsBusy = 'preview' | 'save' | null;

export interface YukumoJsTalkInput {
  version: AquesTalkVersion;
  text: string;
  at1: AquesTalk1Config;
  at2: AquesTalk2Config;
  at10: AquesTalk10Config;
}

export function useYukumoJsTalk() {
  const busy = ref<YukumoJsBusy>(null);

  let engine: TalkEngine | null = null;
  let engineKey = '';
  let kanji2koe: AqKanji2Koe | null = null;
  let kanji2koePending: Promise<AqKanji2Koe> | null = null;
  let stopPlaybackFn: (() => void) | null = null;

  function stopPlayback(): void {
    stopPlaybackFn?.();
    stopPlaybackFn = null;
  }

  async function destroyEngine(): Promise<void> {
    const current = engine;
    engine = null;
    engineKey = '';
    if (current) {
      await current.destroy();
    }
  }

  async function ensureEngine(input: YukumoJsTalkInput): Promise<TalkEngine> {
    const key = engineCacheKey(input.version, input.at1, input.at2);
    if (engine != null && engineKey === key) {
      return engine;
    }
    await destroyEngine();
    const next = await loadTalkEngine(input.version, input.at1, input.at2);
    engine = next;
    engineKey = key;
    return next;
  }

  async function ensureKanji2Koe(): Promise<AqKanji2Koe> {
    if (kanji2koe != null) {
      return kanji2koe;
    }
    if (kanji2koePending == null) {
      kanji2koePending = loadKanji2Koe();
    }
    kanji2koe = await kanji2koePending;
    return kanji2koe;
  }

  async function makeWav(input: YukumoJsTalkInput): Promise<Uint8Array> {
    const [talkEngine, converter] = await Promise.all([
      ensureEngine(input),
      ensureKanji2Koe(),
    ]);
    const koe = converter.convert(input.text.trim());
    return synthesize(
      talkEngine,
      input.version,
      koe,
      input.at1,
      input.at2,
      input.at10,
    );
  }

  async function preview(input: YukumoJsTalkInput): Promise<void> {
    if (busy.value != null) return;
    busy.value = 'preview';
    try {
      const wav = await makeWav(input);
      stopPlayback();
      const playback = createPlayback(wav, () => {
        stopPlaybackFn = null;
      });
      stopPlaybackFn = playback.stop;
      await playback.audio.play();
    } catch (error) {
      stopPlayback();
      throw error;
    } finally {
      busy.value = null;
    }
  }

  async function save(input: YukumoJsTalkInput): Promise<void> {
    if (busy.value != null) return;
    busy.value = 'save';
    try {
      const wav = await makeWav(input);
      downloadWav(wav, wavFilename(input.text));
    } finally {
      busy.value = null;
    }
  }

  onBeforeUnmount(() => {
    stopPlayback();
    void destroyEngine();
    const converter = kanji2koe;
    kanji2koe = null;
    kanji2koePending = null;
    void converter?.destroy();
  });

  return {
    busy,
    preview,
    save,
    stopPlayback,
  };
}
