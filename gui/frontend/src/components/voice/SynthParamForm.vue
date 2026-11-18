<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { FileAudio, Trash2 } from '@lucide/vue';
import { Snackbar } from '@varlet/ui';
import EngineBadge from '../profiles/EngineBadge.vue';
import { usePresetCatalog } from '../../backend/presets';
import type { SynthParam } from '../../backend/types';
import { bytesToBase64 } from '../../backend/wav';
import SliderField from './SliderField.vue';

const props = defineProps<{
  modelValue: SynthParam;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: SynthParam];
}>();

const { t } = useI18n();
const { ensure, presetNames, presetLabel } = usePresetCatalog();

const loadingPresets = ref(false);
const phontInput = ref<HTMLInputElement | null>(null);

const presetOptions = computed(() => {
  const names = presetNames(props.modelValue.version);
  if (names.length === 0) {
    return [{ value: props.modelValue.preset, label: `#${props.modelValue.preset}` }];
  }
  const options = names.map((name, index) => ({
    value: index,
    label: name || `#${index}`,
  }));
  if (!options.some((option) => option.value === props.modelValue.preset)) {
    options.push({
      value: props.modelValue.preset,
      label: `#${props.modelValue.preset}`,
    });
  }
  return options;
});

const customPreset = computed(
  () => props.modelValue.version === 10 && props.modelValue.preset === 0,
);

const phontName = computed(() => {
  const encoded = props.modelValue.phontData;
  if (!encoded) return '';
  return t('pages.profiles.wizard.phontLoaded', {
    n: Math.ceil((encoded.length * 3) / 4),
  });
});

function patch(part: Partial<SynthParam>): void {
  emit('update:modelValue', { ...props.modelValue, ...part });
}

function onPresetChange(value: string | number): void {
  patch({ preset: Number(value) });
}

function onBaseChange(value: string | number): void {
  patch({ base: Number(value) });
}

async function loadPresets(version: number): Promise<void> {
  loadingPresets.value = true;
  try {
    await ensure(version);
  } catch (cause) {
    Snackbar({
      content: cause instanceof Error ? cause.message : String(cause),
      type: 'error',
    });
  } finally {
    loadingPresets.value = false;
  }
}

function openPhontPicker(): void {
  phontInput.value?.click();
}

async function onPhontChange(event: Event): Promise<void> {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  input.value = '';
  if (!file) return;
  try {
    const bytes = new Uint8Array(await file.arrayBuffer());
    patch({ phontData: bytesToBase64(bytes) });
  } catch (cause) {
    Snackbar({
      content: cause instanceof Error ? cause.message : String(cause),
      type: 'error',
    });
  }
}

function clearPhont(): void {
  const next = { ...props.modelValue };
  delete next.phontData;
  emit('update:modelValue', next);
}

watch(() => props.modelValue.version, loadPresets, { immediate: true });
</script>

<template>
  <div class="flex flex-col gap-3">
    <div class="flex items-center gap-2">
      <EngineBadge :version="modelValue.version" />
      <span class="text-sm text-on-surface-variant">
        {{ presetLabel(modelValue.version, modelValue.preset) }}
      </span>
      <span
        v-if="loadingPresets"
        class="text-xs text-on-surface-variant"
      >
        …
      </span>
    </div>

    <var-select
      :model-value="modelValue.preset"
      size="small"
      variant="outlined"
      :placeholder="t('pages.profiles.wizard.preset')"
      @update:model-value="onPresetChange"
    >
      <var-option
        v-for="option in presetOptions"
        :key="option.value"
        :value="option.value"
        :label="option.label"
      />
    </var-select>

    <SliderField
      :model-value="modelValue.speed"
      :label="t('pages.profiles.wizard.speed')"
      :min="50"
      :max="300"
      @update:model-value="patch({ speed: $event })"
    />

    <div
      v-if="modelValue.version === 2"
      class="flex flex-col gap-1.5"
    >
      <p class="text-sm font-medium text-text">
        {{ t('pages.profiles.wizard.phont') }}
      </p>
      <div class="flex items-center gap-2">
        <var-button
          size="small"
          outlined
          @click="openPhontPicker"
        >
          <FileAudio
            :size="16"
            aria-hidden="true"
            class="mr-1.5"
          />
          {{ t('pages.profiles.wizard.phontChoose') }}
        </var-button>
        <var-button
          v-if="modelValue.phontData"
          size="small"
          text
          @click="clearPhont"
        >
          <Trash2
            :size="16"
            aria-hidden="true"
          />
        </var-button>
        <span
          v-if="phontName"
          class="truncate text-sm text-on-surface-variant"
        >
          {{ phontName }}
        </span>
      </div>
      <input
        ref="phontInput"
        class="hidden"
        type="file"
        accept=".phont"
        @change="onPhontChange"
      />
    </div>

    <template v-if="customPreset">
      <var-select
        :model-value="modelValue.base"
        size="small"
        variant="outlined"
        :placeholder="t('pages.profiles.wizard.bas')"
        @update:model-value="onBaseChange"
      >
        <var-option
          :value="0"
          :label="t('pages.profiles.wizard.basF1E')"
        />
        <var-option
          :value="1"
          :label="t('pages.profiles.wizard.basF2E')"
        />
        <var-option
          :value="2"
          :label="t('pages.profiles.wizard.basM1E')"
        />
      </var-select>
      <SliderField
        :model-value="modelValue.volume"
        :label="t('pages.profiles.wizard.vol')"
        :min="0"
        :max="300"
        @update:model-value="patch({ volume: $event })"
      />
      <SliderField
        :model-value="modelValue.pitch"
        :label="t('pages.profiles.wizard.pit')"
        :min="20"
        :max="200"
        @update:model-value="patch({ pitch: $event })"
      />
      <SliderField
        :model-value="modelValue.accent"
        :label="t('pages.profiles.wizard.acc')"
        :min="0"
        :max="200"
        @update:model-value="patch({ accent: $event })"
      />
      <SliderField
        :model-value="modelValue.lmd"
        :label="t('pages.profiles.wizard.lmd')"
        :min="0"
        :max="200"
        @update:model-value="patch({ lmd: $event })"
      />
      <SliderField
        :model-value="modelValue.fsc"
        :label="t('pages.profiles.wizard.fsc')"
        :min="50"
        :max="200"
        @update:model-value="patch({ fsc: $event })"
      />
    </template>
  </div>
</template>
