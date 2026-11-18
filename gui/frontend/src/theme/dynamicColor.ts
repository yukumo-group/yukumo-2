import { Hct, hexFromArgb } from '@material/material-color-utilities';

/** HCT hue matching Material MD3 purple (#6750A4). */
export const DEFAULT_THEME_HUE = 300;

const SOURCE_CHROMA = 48;
const SOURCE_TONE = 40;

export function normalizeHue(hue: number): number {
  if (!Number.isFinite(hue)) {
    return DEFAULT_THEME_HUE;
  }
  const wrapped = ((Math.round(hue) % 360) + 360) % 360;
  return wrapped;
}

/** Source color passed to `m3e-theme`. The component derives the scheme. */
export function sourceHexForHue(hue: number): string {
  const normalized = normalizeHue(hue);
  return hexFromArgb(Hct.from(normalized, SOURCE_CHROMA, SOURCE_TONE).toInt());
}
