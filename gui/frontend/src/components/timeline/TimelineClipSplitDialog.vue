<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  show: boolean;
  text: string;
  initialIndex: number;
}>();

const emit = defineEmits<{
  'update:show': [value: boolean];
  confirm: [index: number];
}>();

const { t } = useI18n();

const maxIndex = computed(() => Math.max(1, props.text.length - 1));

const index = ref(clampIndex(props.initialIndex));

function clampIndex(value: number): number {
  if (props.text.length <= 1) return 1;
  return Math.min(maxIndex.value, Math.max(1, Math.round(value)));
}

watch(
  () => props.show,
  (show) => {
    if (show) index.value = clampIndex(props.initialIndex);
  },
);

const leftText = computed(() => props.text.slice(0, index.value));
const rightText = computed(() => props.text.slice(index.value));
const canSplit = computed(() => props.text.length > 1);

function setShow(value: boolean): void {
  emit('update:show', value);
}

function onSlider(event: Event): void {
  const thumb = event.target as HTMLElement & { value?: number };
  if (typeof thumb.value !== 'number') return;
  index.value = clampIndex(thumb.value);
}

function onConfirm(): void {
  if (!canSplit.value) return;
  emit('confirm', index.value);
  setShow(false);
}
</script>

<template>
  <m3e-dialog
    class="timeline-clip-split-dialog"
    .open="show"
    @closed="setShow(false)"
  >
    <span slot="header">{{ t('pages.generate.timeline.splitDialogTitle') }}</span>
    <div class="flex flex-col gap-4 text-left">
      <p class="text-sm text-on-surface-variant">
        {{ t('pages.generate.timeline.splitDialogHint') }}
      </p>
      <div
        class="max-h-40 overflow-auto whitespace-pre-wrap break-all rounded-lg bg-surface-container px-3 py-2 text-sm text-on-surface"
        :aria-label="t('pages.generate.timeline.splitPreviewAriaLabel')"
      >
        <span>{{ leftText }}</span>
        <span
          class="mx-0.5 inline-block font-bold text-primary"
          aria-hidden="true"
        >|</span>
        <span>{{ rightText }}</span>
      </div>
      <m3e-slider
        .min="1"
        .max="maxIndex"
        .step="1"
        labelled
        .disabled="!canSplit"
        @input="onSlider"
      >
        <m3e-slider-thumb .value="index" />
      </m3e-slider>
    </div>
    <div slot="actions" class="flex w-full justify-end gap-2">
      <m3e-button variant="text" @click="setShow(false)">
        {{ t('pages.generate.timeline.dialogCancel') }}
      </m3e-button>
      <m3e-button variant="filled" .disabled="!canSplit" @click="onConfirm">
        {{ t('pages.generate.timeline.splitConfirm') }}
      </m3e-button>
    </div>
  </m3e-dialog>
</template>

<style>
.timeline-clip-split-dialog {
  --m3e-dialog-max-width: min(32rem, 92vw);
}
</style>
