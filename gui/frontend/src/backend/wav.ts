export function base64ToBytes(encoded: string): Uint8Array {
  if (!encoded) return new Uint8Array();
  const binary = atob(encoded);
  const out = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i += 1) {
    out[i] = binary.charCodeAt(i);
  }
  return out;
}

export function bytesToBase64(bytes: Uint8Array): string {
  let binary = '';
  const chunkSize = 0x8000;
  for (let i = 0; i < bytes.length; i += chunkSize) {
    binary += String.fromCharCode(...bytes.subarray(i, i + chunkSize));
  }
  return btoa(binary);
}

export function wavBlob(wav: Uint8Array): Blob {
  return new Blob([new Uint8Array(wav)], { type: 'audio/wav' });
}

export function wavFilename(content: string): string {
  const slug = content
    .replace(/[\\/:*?"<>|\r\n]+/g, '')
    .replace(/\s+/g, '_')
    .replace(/^\.+/, '')
    .slice(0, 80)
    .replace(/\.+$/, '');
  return `yukumo_${slug || 'audio'}.wav`;
}

export function downloadWav(wav: Uint8Array, filename = 'yukumo.wav'): void {
  const url = URL.createObjectURL(wavBlob(wav));
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  link.click();
  URL.revokeObjectURL(url);
}

export function createPlayback(
  wav: Uint8Array,
  onEnded: () => void,
): {
  audio: HTMLAudioElement;
  stop: () => void;
} {
  const url = URL.createObjectURL(wavBlob(wav));
  const audio = new Audio(url);
  let released = false;

  const release = () => {
    if (released) return;
    released = true;
    audio.pause();
    audio.currentTime = 0;
    URL.revokeObjectURL(url);
  };

  audio.addEventListener('ended', () => {
    release();
    onEnded();
  });

  return { audio, stop: release };
}
