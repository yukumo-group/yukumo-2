import { onMounted, onUnmounted, ref, watch, type Ref } from 'vue';
import type { M3eSplitPaneElement, SplitPaneOrientation } from '@m3e/web/split-pane';
import { SPLIT_AUTO_HIDE_THRESHOLD_PX, clampSplitSize } from './splitSize';

const SPLIT_PANE_STYLE = `
.start,
.end {
  min-width: 0 !important;
  min-height: 0 !important;
}
.start {
  flex-grow: 0 !important;
  flex-shrink: 0 !important;
}
:host(.split-collapsed[data-collapse='end']) .start,
:host(.split-collapsed[data-collapse='start']) .end {
  flex: 1 1 auto !important;
}
:host(.split-collapsed[data-collapse='end']) .end,
:host(.split-collapsed[data-collapse='start']) .start,
:host(.split-collapsed) .drag-handle {
  flex: 0 0 0px !important;
  width: 0 !important;
  height: 0 !important;
  min-width: 0 !important;
  min-height: 0 !important;
  overflow: hidden !important;
  padding: 0 !important;
}
`;

export interface PixelSplitOptions {
  open: Ref<boolean>;
  size: Ref<number>;
  /** `start` sizes the first slot. `end` sizes the second slot. */
  side: 'start' | 'end';
  orientation: Exclude<SplitPaneOrientation, 'auto'>;
  minPx: number;
  maxPx: number;
  reservePx: number;
}

export function useM3eSplit(options: PixelSplitOptions): {
  splitRef: Ref<M3eSplitPaneElement | null>;
  onSplitInput: (event: Event) => void;
  onSplitChange: (event: Event) => void;
} {
  const splitRef = ref<M3eSplitPaneElement | null>(null);
  let fromUser = false;
  let observer: ResizeObserver | undefined;

  function span(el: M3eSplitPaneElement): number {
    return options.orientation === 'vertical' ? el.clientHeight : el.clientWidth;
  }

  function valueFor(px: number, container: number): number {
    const percent = container > 0 ? (px / container) * 100 : 0;
    return options.side === 'start' ? percent : 100 - percent;
  }

  function pxFor(value: number, container: number): number {
    const startPx = (value / 100) * container;
    return options.side === 'start' ? startPx : container - startPx;
  }

  function ensurePaneStyle(el: M3eSplitPaneElement): void {
    const root = el.shadowRoot;
    if (!root) return;
    const existing = root.getElementById('yukumo-split-pane');
    if (existing?.textContent === SPLIT_PANE_STYLE) return;
    const style = existing ?? document.createElement('style');
    style.id = 'yukumo-split-pane';
    style.textContent = SPLIT_PANE_STYLE;
    if (!existing) root.append(style);
  }

  function apply(animate: boolean): void {
    const el = splitRef.value;
    if (!el) return;
    ensurePaneStyle(el);
    el.dataset.collapse = options.side;
    const container = span(el);
    if (container <= 0) return;

    const openPx = clampSplitSize(
      options.size.value,
      options.minPx,
      options.maxPx,
      container,
      options.reservePx,
    );
    const maxPane = Math.min(
      options.maxPx,
      Math.max(options.minPx, container - options.reservePx),
    );
    const minPercent = (options.minPx / container) * 100;
    const maxPercent = (maxPane / container) * 100;

    if (options.open.value) {
      if (options.side === 'start') {
        el.min = minPercent;
        el.max = maxPercent;
      } else {
        el.min = 100 - maxPercent;
        el.max = 100 - minPercent;
      }
      el.disabled = false;
      const next = valueFor(openPx, container);
      if (animate) el.snapToValue(next);
      else el.value = next;
      return;
    }

    el.min = options.side === 'start' ? 0 : 100 - maxPercent;
    el.max = 100;
    el.disabled = true;
    const closed = options.side === 'start' ? 0 : 100;
    if (animate) el.snapToValue(closed);
    else el.value = closed;
  }

  function onSplit(event: Event): void {
    const el = event.currentTarget;
    if (!(el instanceof HTMLElement) || !('value' in el)) return;
    const split = el as M3eSplitPaneElement;
    if (!options.open.value) return;
    const container = span(split);
    if (container <= 0) return;

    const px = pxFor(split.value, container);
    if (px < options.minPx - SPLIT_AUTO_HIDE_THRESHOLD_PX) {
      fromUser = true;
      options.open.value = false;
      fromUser = false;
      apply(true);
      return;
    }

    options.size.value = clampSplitSize(
      px,
      options.minPx,
      options.maxPx,
      container,
      options.reservePx,
    );
  }

  function observe(el: M3eSplitPaneElement): void {
    observer?.disconnect();
    observer = new ResizeObserver(() => {
      if (fromUser) return;
      if (splitRef.value?.matches(':state(--pressed)')) return;
      apply(false);
    });
    observer.observe(el);
  }

  function sync(animate: boolean): void {
    apply(animate);
    const el = splitRef.value;
    if (!el) return;
    if (span(el) <= 0) {
      requestAnimationFrame(() => apply(false));
    }
  }

  watch(options.open, () => {
    if (fromUser) return;
    sync(true);
  });

  watch(splitRef, (el) => {
    if (!el) return;
    sync(false);
    observe(el);
  });

  onMounted(() => {
    sync(false);
    const el = splitRef.value;
    if (el) observe(el);
  });

  onUnmounted(() => {
    observer?.disconnect();
  });

  return {
    splitRef,
    onSplitInput: onSplit,
    onSplitChange: onSplit,
  };
}
