<script setup lang="ts">
import { computed } from 'vue';
import { useId } from 'vue';
import { useI18n } from 'vue-i18n';
import { Plus, Trash2 } from '@lucide/vue';
import EngineBadge from '../profiles/EngineBadge.vue';
import {
  createDefaultParam,
  normalizeVoice,
  removeParam,
  upsertParam,
  voiceVersions,
} from '../../backend/params';
import type { AquesTalkVersion, Voice } from '../../backend/types';
import { AQUES_TALK_VERSIONS } from '../../backend/types';
import SynthParamForm from './SynthParamForm.vue';

const props = defineProps<{
  modelValue: Voice;
  removable?: boolean;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: Voice];
  remove: [];
}>();

const { t } = useI18n();
const nameId = useId();

const missingVersions = computed(() => {
  const present = new Set(voiceVersions(props.modelValue));
  return AQUES_TALK_VERSIONS.filter((version) => !present.has(version));
});

function update(next: Voice): void {
  emit('update:modelValue', next);
}

function onNameInput(event: Event): void {
  const input = event.target;
  if (!(input instanceof HTMLInputElement)) return;
  update({ ...props.modelValue, name: input.value });
}

function onParamChange(
  version: number,
  value: Voice['params'][number],
): void {
  update(upsertParam(props.modelValue, { ...value, version }));
}

function addParam(version: AquesTalkVersion): void {
  update(upsertParam(props.modelValue, createDefaultParam(version)));
}

function removeVersion(version: number): void {
  update(removeParam(props.modelValue, version));
}

function removeVoice(): void {
  emit('remove');
}

function onBlur(): void {
  update(normalizeVoice(props.modelValue));
}
</script>

<template>
  <div class="flex flex-col gap-3 rounded-xl border border-outline/30 p-3">
    <div class="flex items-center gap-2">
      <m3e-form-field class="min-w-0 flex-1">
        <label slot="label" :for="nameId">
          {{ t('pages.profiles.wizard.voiceName') }}
        </label>
        <input
          :id="nameId"
          :value="modelValue.name"
          :aria-label="t('pages.profiles.wizard.voiceName')"
          @input="onNameInput"
          @blur="onBlur"
        />
      </m3e-form-field>
      <m3e-icon-button
        v-if="removable"
        class="tone-error"
        :aria-label="t('pages.profiles.wizard.removeVoice')"
        @click="removeVoice"
      >
        <Trash2
          :size="16"
          aria-hidden="true"
        />
      </m3e-icon-button>
    </div>

    <div
      v-for="param in modelValue.params"
      :key="param.version"
      class="flex flex-col gap-2 rounded-lg bg-surface-container-low p-3"
    >
      <div class="flex items-center justify-between gap-2">
        <span class="flex items-center gap-2 text-sm font-medium text-on-surface">
          <EngineBadge :version="param.version" />
        </span>
        <m3e-icon-button
          v-if="modelValue.params.length > 1"
          :aria-label="t('pages.profiles.wizard.removeParam')"
          @click="removeVersion(param.version)"
        >
          <Trash2
            :size="14"
            aria-hidden="true"
          />
        </m3e-icon-button>
      </div>
      <SynthParamForm
        :model-value="param"
        @update:model-value="onParamChange(param.version, $event)"
      />
    </div>

    <div
      v-if="missingVersions.length > 0"
      class="flex flex-wrap items-center gap-2"
    >
      <span class="text-sm text-on-surface-variant">
        {{ t('pages.profiles.wizard.addParam') }}
      </span>
      <m3e-button
        v-for="version in missingVersions"
        :key="version"
        size="small"
        variant="outlined"
        @click="addParam(version)"
      >
        <Plus slot="icon" :size="14" aria-hidden="true" />
        <EngineBadge :version="version" />
      </m3e-button>
    </div>
  </div>
</template>
