<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useId } from 'vue';
import { useI18n } from 'vue-i18n';
import { FileAudio, Trash2 } from '@lucide/vue';
import type { M3eSelectElement } from '@m3e/web/select';
import { notify } from '../../ui/notify';
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
const presetSelectId = useId();
const baseSelectId = useId();
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

function selectedNumber(event: Event): number | null {
  const select = event.currentTarget as M3eSelectElement;
  if (typeof select.value !== 'string' || select.value === '') return null;
  const next = Number(select.value);
  return Number.isFinite(next) ? next : null;
}

function onPresetChange(event: Event): void {
  const next = selectedNumber(event);
  if (next === null) return;
  patch({ preset: next });
}

function onBaseChange(event: Event): void {
  const next = selectedNumber(event);
  if (next === null) return;
  patch({ base: next });
}

async function loadPresets(version: number): Promise<void> {
  loadingPresets.value = true;
  try {
    await ensure(version);
  } catch (cause) {
    notify(cause instanceof Error ? cause.message : String(cause));
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
    notify(cause instanceof Error ? cause.message : String(cause));
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

    <m3e-form-field>
      <label slot="label" :for="presetSelectId">
        {{ t('pages.profiles.wizard.preset') }}
      </label>
      <m3e-select
        :id="presetSelectId"
        .value="String(modelValue.preset)"
        @change="onPresetChange"
      >
        <m3e-option
          v-for="option in presetOptions"
          :key="option.value"
          :value="String(option.value)"
        >
          {{ option.label }}
        </m3e-option>
      </m3e-select>
    </m3e-form-field>

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
      <p class="text-sm font-medium text-on-surface">
        {{ t('pages.profiles.wizard.phont') }}
      </p>
      <div class="flex items-center gap-2">
        <m3e-button
          size="small"
          variant="outlined"
          @click="openPhontPicker"
        >
          <FileAudio slot="icon" :size="16" aria-hidden="true" />
          {{ t('pages.profiles.wizard.phontChoose') }}
        </m3e-button>
        <m3e-icon-button
          v-if="modelValue.phontData"
          :aria-label="t('pages.profiles.wizard.phontChoose')"
          @click="clearPhont"
        >
          <Trash2
            :size="16"
            aria-hidden="true"
          />
        </m3e-icon-button>
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
      <m3e-form-field>
        <label slot="label" :for="baseSelectId">
          {{ t('pages.profiles.wizard.bas') }}
        </label>
        <m3e-select
          :id="baseSelectId"
          .value="String(modelValue.base)"
          @change="onBaseChange"
        >
          <m3e-option value="0">
            {{ t('pages.profiles.wizard.basF1E') }}
          </m3e-option>
          <m3e-option value="1">
            {{ t('pages.profiles.wizard.basF2E') }}
          </m3e-option>
          <m3e-option value="2">
            {{ t('pages.profiles.wizard.basM1E') }}
          </m3e-option>
        </m3e-select>
      </m3e-form-field>
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
