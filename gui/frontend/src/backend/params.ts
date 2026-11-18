import type { AquesTalkVersion, SynthParam, Voice } from './types';
import { AQUES_TALK_VERSIONS, isAquesTalkVersion } from './types';

const SPEED_MIN = 50;
const SPEED_MAX = 300;

export function createDefaultParam(version: AquesTalkVersion): SynthParam {
  if (version === 10) {
    return {
      version,
      preset: 0,
      speed: 100,
      base: 0,
      volume: 100,
      pitch: 100,
      accent: 100,
      lmd: 100,
      fsc: 100,
    };
  }
  return {
    version,
    preset: 0,
    speed: 100,
    base: 0,
    volume: 0,
    pitch: 0,
    accent: 0,
    lmd: 0,
    fsc: 0,
  };
}

export function createDefaultVoice(
  name = 'default',
  versions: readonly AquesTalkVersion[] = [1],
): Voice {
  return {
    name,
    params: versions.map((version) => createDefaultParam(version)),
  };
}

function clamp(value: number, min: number, max: number): number {
  const rounded = Math.round(Number.isFinite(value) ? value : 0);
  return Math.min(max, Math.max(min, rounded));
}

export function normalizeParam(param: SynthParam): SynthParam {
  const version = isAquesTalkVersion(param.version) ? param.version : 1;
  const out: SynthParam = {
    ...createDefaultParam(version),
    ...param,
    version,
    preset: Math.max(0, Math.trunc(param.preset || 0)),
    speed: clamp(param.speed, SPEED_MIN, SPEED_MAX),
  };
  if (version !== 2) {
    delete out.phontData;
  }
  if (version !== 10) {
    out.base = 0;
    out.volume = 0;
    out.pitch = 0;
    out.accent = 0;
    out.lmd = 0;
    out.fsc = 0;
    return out;
  }
  out.base = clamp(param.base, 0, 2);
  out.volume = clamp(param.volume, 0, 300);
  out.pitch = clamp(param.pitch, 20, 200);
  out.accent = clamp(param.accent, 0, 200);
  out.lmd = clamp(param.lmd, 0, 200);
  out.fsc = clamp(param.fsc, 50, 200);
  return out;
}

function parameterRank(param: SynthParam): number {
  const index = AQUES_TALK_VERSIONS.indexOf(param.version as AquesTalkVersion);
  return index < 0 ? AQUES_TALK_VERSIONS.length : index;
}

export function normalizeVoice(voice: Voice): Voice {
  const byVersion = new Map<number, SynthParam>();
  for (const param of voice.params ?? []) {
    if (!isAquesTalkVersion(param.version)) continue;
    byVersion.set(param.version, normalizeParam(param));
  }
  return {
    name: voice.name.trim(),
    params: [...byVersion.values()].sort(
      (a, b) => parameterRank(a) - parameterRank(b),
    ),
  };
}

export function paramForVersion(
  voice: Voice,
  version: number,
): SynthParam | undefined {
  return voice.params.find((param) => param.version === version);
}

export function upsertParam(voice: Voice, param: SynthParam): Voice {
  const next = voice.params.filter((item) => item.version !== param.version);
  next.push(param);
  return normalizeVoice({ ...voice, params: next });
}

export function removeParam(voice: Voice, version: number): Voice {
  return normalizeVoice({
    ...voice,
    params: voice.params.filter((param) => param.version !== version),
  });
}

export function voiceSupportsVersion(voice: Voice, version: number): boolean {
  return voice.params.some((param) => param.version === version);
}

export function voiceVersions(voice: Voice): AquesTalkVersion[] {
  return AQUES_TALK_VERSIONS.filter((version) =>
    voiceSupportsVersion(voice, version),
  );
}
