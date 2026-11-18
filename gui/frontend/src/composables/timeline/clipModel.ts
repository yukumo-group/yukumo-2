import {
  createDefaultVoice,
  voiceSupportsVersion,
} from '../../backend/params';
import type { AquesTalkVersion } from '../../backend/types';
import {
  TIMELINE_CLIP_MIN_DURATION_SEC,
  type TimelineClip,
  type TimelineClipInsert,
  type TimelineClipSpeaker,
} from '../../types/timeline';

export function createDefaultSpeaker(): TimelineClipSpeaker {
  return {
    characterName: null,
    voice: createDefaultVoice('custom', [1, 2, 10]),
  };
}

export function cloneSpeaker(
  speaker: TimelineClipSpeaker,
): TimelineClipSpeaker {
  return {
    characterName: speaker.characterName,
    voice: {
      name: speaker.voice.name,
      params: speaker.voice.params.map((param) => ({ ...param })),
    },
  };
}

export function speakerSupportsEngine(
  speaker: TimelineClipSpeaker,
  engine: AquesTalkVersion,
): boolean {
  return voiceSupportsVersion(speaker.voice, engine);
}

export function clipPayload(clip: TimelineClip): TimelineClipInsert {
  return {
    trackId: clip.trackId,
    startSec: clip.startSec,
    durationSec: clip.durationSec,
    text: clip.text,
    speaker: cloneSpeaker(clip.speaker),
    volume: clip.volume,
    pan: clip.pan,
    muted: clip.muted,
    color: clip.color,
  };
}

export function createClipDefaults(
  trackId: string,
  startSec: number,
  durationSec: number,
  speaker?: TimelineClipSpeaker,
): TimelineClipInsert {
  return {
    trackId,
    startSec: Math.max(0, startSec),
    durationSec,
    text: '',
    speaker: speaker ? cloneSpeaker(speaker) : createDefaultSpeaker(),
    volume: 0.75,
    pan: 0,
    muted: false,
    color: 'primary',
  };
}

function formatSigned(n: number, digits: number): string {
  const rounded = Number(n.toFixed(digits));
  if (rounded === 0) return `+${(0).toFixed(digits)}`;
  if (rounded > 0) return `+${rounded.toFixed(digits)}`;
  return rounded.toFixed(digits);
}

export function formatVolumeDb(volume: number): string {
  if (volume <= 0) return '-∞ dB';
  return `${formatSigned(20 * Math.log10(volume), 1)} dB`;
}

export function formatPanPercent(pan: number): string {
  return `${Math.round(pan * 100)}%`;
}

export function formatClipLevels(volume: number, pan: number): string {
  return `${formatVolumeDb(volume)}  ${formatPanPercent(pan)}`;
}

export function isValidClipTimeSplit(
  startSec: number,
  durationSec: number,
  atSec: number,
): boolean {
  const left = atSec - startSec;
  const right = startSec + durationSec - atSec;
  return (
    left >= TIMELINE_CLIP_MIN_DURATION_SEC &&
    right >= TIMELINE_CLIP_MIN_DURATION_SEC
  );
}

export function splitIndexFromTime(
  startSec: number,
  durationSec: number,
  textLength: number,
  atSec: number,
): number {
  if (textLength <= 1) return 1;
  const ratio = (atSec - startSec) / durationSec;
  return Math.min(textLength - 1, Math.max(1, Math.round(ratio * textLength)));
}
