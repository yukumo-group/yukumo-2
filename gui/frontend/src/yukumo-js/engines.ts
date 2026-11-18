import {
  AquesTalk10,
  load,
  loadAquesTalk2,
  loadAquesTalk10,
  type AquesTalk2Voice,
  type AquesTalk10Voice,
  type AqtkVoice,
  type Voice,
} from 'yukumo.js';
import { loadAqKanji2Koe, type AqKanji2Koe } from 'yukumo.js/lang/kanji2koe';
import type {
  AquesTalk1Config,
  AquesTalk10Config,
  AquesTalk2Config,
  AquesTalkVersion,
} from '../types/profile';
import { generateDevKey } from './kg';

export type TalkEngine = Awaited<ReturnType<typeof load>> | Awaited<
  ReturnType<typeof loadAquesTalk2>
> | AquesTalk10;

export const MEMORY_SIZE = 1024 * 1024 * 1024;

const AQ1_VOICES: readonly Voice[] = [
  'dvd',
  'f1',
  'f2',
  'imd1',
  'jgr',
  'm1',
  'm2',
  'r1',
];

const AQ2_VOICES: readonly AquesTalk2Voice[] = [
  'f4',
  'rm3',
  'm5',
  'mf2',
  'mf1',
  'rm',
  'f1c',
  'm4b',
  'robo',
  'rb3',
  'rb2',
  'f3a',
  'huskey',
  'yukkuri',
  'defo1',
  'teto1',
  'momo1',
];

function isAq1Voice(value: string): value is Voice {
  return (AQ1_VOICES as readonly string[]).includes(value);
}

function isAq2Voice(value: string): value is AquesTalk2Voice {
  return (AQ2_VOICES as readonly string[]).includes(value);
}

export function resolveAt1Voice(voiceId: string): Voice {
  const id = voiceId.trim().toLowerCase();
  if (isAq1Voice(id)) return id;
  throw new Error(`Unsupported AquesTalk1 voice '${voiceId}'`);
}

export function resolveAt2Voice(phontName: string): AquesTalk2Voice {
  const key = phontName.trim().toLowerCase().replace(/^(aq_|ar_)/, '');
  if (isAq2Voice(key)) return key;
  throw new Error(`Unsupported AquesTalk2 phont '${phontName}'`);
}

export function at10Params(config: AquesTalk10Config): AqtkVoice {
  return {
    bas: config.bas,
    spd: config.spd,
    vol: config.vol,
    pit: config.pit,
    acc: config.acc,
    lmd: config.lmd,
    fsc: config.fsc,
  };
}

export function engineCacheKey(
  version: AquesTalkVersion,
  at1: AquesTalk1Config,
  at2: AquesTalk2Config,
): string {
  if (version === 1) return `1:${resolveAt1Voice(at1.voiceId)}`;
  if (version === 2) return `2:${resolveAt2Voice(at2.phontName)}`;
  return '10';
}

export async function loadTalkEngine(
  version: AquesTalkVersion,
  at1: AquesTalk1Config,
  at2: AquesTalk2Config,
): Promise<TalkEngine> {
  const options = { memorySize: MEMORY_SIZE };
  if (version === 1) {
    const aq1 = await load(resolveAt1Voice(at1.voiceId), options);
    const keyResult = aq1.SetDevKey(generateDevKey({ module: 'aqtk1' }));
    if (keyResult !== 0) {
      throw new Error(`SetDevKey failed with code ${keyResult}`);
    }
    return aq1;
  }
  if (version === 2) {
    return loadAquesTalk2(resolveAt2Voice(at2.phontName), options);
  }
  const aq10 = await loadAquesTalk10('f1' satisfies AquesTalk10Voice, options);
  const keyResult = aq10.SetDevKey(generateDevKey({ module: 'tk10' }));
  if (keyResult !== 0) {
    throw new Error(`SetDevKey failed with code ${keyResult}`);
  }
  return aq10;
}

export async function loadKanji2Koe(): Promise<AqKanji2Koe> {
  const k2k = await loadAqKanji2Koe({ memorySize: MEMORY_SIZE });
  k2k.setDevKey(generateDevKey({ module: 'k2k' }));
  return k2k;
}

export function synthesize(
  engine: TalkEngine,
  version: AquesTalkVersion,
  koe: string,
  at1: AquesTalk1Config,
  at2: AquesTalk2Config,
  at10: AquesTalk10Config,
): Uint8Array {
  if (version === 10) {
    const param = at10Params(at10);
    return (engine as AquesTalk10).run(koe, param.spd, param);
  }
  const speed = version === 1 ? at1.speed : at2.speed;
  return engine.run(koe, speed);
}
