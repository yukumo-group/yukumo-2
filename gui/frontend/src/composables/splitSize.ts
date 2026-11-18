/** Drag this far past min before a hideable panel auto-collapses. */
export const SPLIT_AUTO_HIDE_THRESHOLD_PX = 64;

export function clampSplitSize(
  px: number,
  min: number,
  max: number,
  containerPx?: number,
  reservePx = 0,
): number {
  let resolvedMax = max;
  if (containerPx != null && Number.isFinite(containerPx)) {
    resolvedMax = Math.min(max, Math.max(min, containerPx - reservePx));
  }
  return Math.min(resolvedMax, Math.max(min, Math.round(px)));
}
