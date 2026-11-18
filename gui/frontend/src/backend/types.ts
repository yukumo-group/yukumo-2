import type {
  Character,
  Config,
  SynthParam,
  Voice,
} from '../../bindings/yukumo-script/core/api';

export type { Character, Config, SynthParam, Voice };

export const AQUES_TALK_VERSIONS = [1, 2, 10] as const;

export type AquesTalkVersion = (typeof AQUES_TALK_VERSIONS)[number];

export function isAquesTalkVersion(value: number): value is AquesTalkVersion {
  return (
    value === AQUES_TALK_VERSIONS[0] ||
    value === AQUES_TALK_VERSIONS[1] ||
    value === AQUES_TALK_VERSIONS[2]
  );
}
