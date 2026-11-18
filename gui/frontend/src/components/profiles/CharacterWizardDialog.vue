<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useId } from 'vue';
import { useI18n } from 'vue-i18n';
import { Plus } from '@lucide/vue';
import type { M3eStepperElement } from '@m3e/web/stepper';
import { confirmDialog } from '../../ui/confirm';
import VoiceEditor from '../voice/VoiceEditor.vue';
import {
  createDefaultVoice,
  normalizeVoice,
} from '../../backend/params';
import type { Character, Voice } from '../../backend/types';
import { useCharacters } from '../../store/useCharacters';

type WizardStepKind = 'basics' | 'voices';

const props = defineProps<{
  show: boolean;
  character: Character | null;
}>();

const emit = defineEmits<{
  'update:show': [value: boolean];
}>();

const { t } = useI18n();
const { addCharacter, updateCharacter, deleteCharacter } = useCharacters();
const nameId = useId();
const descriptionId = useId();

const stepIndex = ref(0);
const error = ref('');
const saving = ref(false);
const name = ref('');
const description = ref('');
const voices = ref<Voice[]>([]);

const isEdit = computed(() => props.character !== null);
const stepKinds = computed<WizardStepKind[]>(() => ['basics', 'voices']);
const currentKind = computed(
  () => stepKinds.value[stepIndex.value] ?? 'basics',
);
const isLast = computed(
  () => stepIndex.value >= stepKinds.value.length - 1,
);

const stepLabelKey: Record<WizardStepKind, string> = {
  basics: 'pages.profiles.wizard.steps.basics',
  voices: 'pages.profiles.wizard.steps.voice',
};

function cloneVoice(voice: Voice): Voice {
  return {
    name: voice.name,
    params: (voice.params ?? []).map((param) => ({ ...param })),
  };
}

function hydrate(): void {
  error.value = '';
  stepIndex.value = 0;
  if (!props.character) {
    name.value = '';
    description.value = '';
    voices.value = [createDefaultVoice('default', [1])];
    return;
  }
  name.value = props.character.name;
  description.value = props.character.description;
  voices.value = props.character.voices.map(cloneVoice);
}

watch(
  () => props.show,
  (open) => {
    if (open) hydrate();
  },
);

function setShow(value: boolean): void {
  emit('update:show', value);
}

function cancel(): void {
  setShow(false);
}

function goTo(index: number): void {
  stepIndex.value = Math.max(0, Math.min(index, stepKinds.value.length - 1));
}

function onStepChange(event: Event): void {
  const stepper = event.currentTarget as M3eStepperElement;
  const next = stepper.selectedIndex;
  if (next < 0) return;
  onClickStep(next);
  if (stepIndex.value !== next) {
    stepper.moveTo(stepIndex.value);
  }
}

function validateKind(kind: WizardStepKind | undefined): boolean {
  error.value = '';
  if (kind === 'basics') {
    if (!name.value.trim()) {
      error.value = t('pages.profiles.wizard.nameRequired');
      return false;
    }
    return true;
  }
  if (voices.value.length === 0) {
    error.value = t('pages.profiles.wizard.voicesRequired');
    return false;
  }
  const seen = new Set<string>();
  for (const voice of voices.value) {
    const voiceName = voice.name.trim();
    if (!voiceName) {
      error.value = t('pages.profiles.wizard.voiceNameRequired');
      return false;
    }
    if (seen.has(voiceName)) {
      error.value = t('pages.profiles.wizard.duplicateVoice', {
        name: voiceName,
      });
      return false;
    }
    seen.add(voiceName);
    if (voice.params.length === 0) {
      error.value = t('pages.profiles.wizard.paramsRequired', {
        name: voiceName,
      });
      return false;
    }
  }
  return true;
}

function validateAll(): boolean {
  return stepKinds.value.every((kind) => validateKind(kind));
}

function onClickStep(index: number): void {
  if (index <= stepIndex.value) {
    goTo(index);
    return;
  }
  for (let i = stepIndex.value; i < index; i += 1) {
    if (!validateKind(stepKinds.value[i])) return;
  }
  goTo(index);
}

function goBack(): void {
  error.value = '';
  goTo(stepIndex.value - 1);
}

function onPrimaryAction(): void {
  if (!isEdit.value && !isLast.value) {
    if (!validateKind(currentKind.value)) return;
    goTo(stepIndex.value + 1);
    return;
  }
  void onSave();
}

function addVoice(): void {
  const index = voices.value.length + 1;
  voices.value = [
    ...voices.value,
    createDefaultVoice(index === 1 ? 'default' : `voice-${index}`, [1]),
  ];
}

function removeVoice(index: number): void {
  voices.value = voices.value.filter((_, i) => i !== index);
}

function updateVoice(index: number, voice: Voice): void {
  voices.value = voices.value.map((item, i) => (i === index ? voice : item));
}

async function onSave(): Promise<void> {
  if (!validateAll()) {
    goTo(stepKinds.value.findIndex((kind) => !validateKind(kind)));
    return;
  }
  saving.value = true;
  error.value = '';
  try {
    const next: Character = {
      name: name.value.trim(),
      description: description.value.trim(),
      voices: voices.value.map(normalizeVoice),
    };
    if (props.character) {
      await updateCharacter(props.character.name, next);
    } else {
      await addCharacter(next);
    }
    setShow(false);
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : String(cause);
  } finally {
    saving.value = false;
  }
}

async function onDelete(): Promise<void> {
  const character = props.character;
  if (!character) return;

  const confirmed = await confirmDialog({
    title: t('pages.profiles.wizard.deleteConfirmTitle'),
    message: t('pages.profiles.wizard.deleteConfirmMessage', {
      name: character.name,
    }),
    confirmText: t('pages.profiles.wizard.deleteConfirm'),
    cancelText: t('pages.profiles.wizard.deleteCancel'),
    danger: true,
  });

  if (!confirmed) return;
  try {
    await deleteCharacter(character.name);
    setShow(false);
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : String(cause);
  }
}

function clearError(): void {
  error.value = '';
}
</script>

<template>
  <m3e-dialog
    class="character-wizard-dialog"
    .open="show"
    @closed="setShow(false)"
  >
    <span slot="header">
      {{
        isEdit
          ? t('pages.profiles.wizard.editTitle')
          : t('pages.profiles.wizard.createTitle')
      }}
    </span>
    <m3e-stepper class="text-left" @change="onStepChange">
      <m3e-step
        v-for="(kind, index) in stepKinds"
        :key="kind"
        :for="`wizard-${kind}`"
        .selected="stepIndex === index"
      >
        {{ t(stepLabelKey[kind]) }}
      </m3e-step>
      <m3e-step-panel id="wizard-basics">
        <m3e-form-field>
          <label slot="label" :for="nameId">
            {{ t('pages.profiles.wizard.namePlaceholder') }}
          </label>
          <input
            :id="nameId"
            v-model="name"
            :aria-label="t('pages.profiles.wizard.namePlaceholder')"
            @input="clearError"
          />
        </m3e-form-field>
        <m3e-form-field>
          <label slot="label" :for="descriptionId">
            {{ t('pages.profiles.wizard.descriptionPlaceholder') }}
          </label>
          <textarea
            :id="descriptionId"
            v-model="description"
            rows="3"
            :aria-label="t('pages.profiles.wizard.descriptionPlaceholder')"
          />
        </m3e-form-field>
      </m3e-step-panel>
      <m3e-step-panel id="wizard-voices">
        <div class="flex items-center justify-between gap-2">
          <p class="text-base font-medium text-on-surface">
            {{ t('pages.profiles.wizard.voices') }}
          </p>
          <m3e-button size="small" variant="outlined" @click="addVoice">
            <Plus slot="icon" :size="14" aria-hidden="true" />
            {{ t('pages.profiles.wizard.addVoice') }}
          </m3e-button>
        </div>
        <div
          v-if="voices.length === 0"
          class="rounded-lg bg-surface-container px-3 py-4 text-sm text-on-surface-variant"
        >
          {{ t('pages.profiles.wizard.noVoices') }}
        </div>
        <VoiceEditor
          v-for="(voice, index) in voices"
          :key="index"
          :model-value="voice"
          removable
          @update:model-value="updateVoice(index, $event)"
          @remove="removeVoice(index)"
        />
      </m3e-step-panel>
    </m3e-stepper>
    <p
      v-if="error"
      class="rounded-lg bg-error/15 px-3 py-2 text-base font-medium text-error"
      role="alert"
    >
      {{ error }}
    </p>
    <div slot="actions" class="flex w-full items-center gap-2">
      <m3e-button
        v-if="isEdit"
        class="tone-error"
        variant="text"
        @click="onDelete"
      >
        {{ t('pages.profiles.wizard.delete') }}
      </m3e-button>
      <div class="ml-auto flex gap-2">
        <m3e-button variant="text" @click="cancel">
          {{ t('pages.profiles.wizard.cancel') }}
        </m3e-button>
        <m3e-button
          v-if="!isEdit"
          variant="text"
          .disabled="stepIndex === 0"
          @click="goBack"
        >
          {{ t('pages.profiles.wizard.back') }}
        </m3e-button>
        <m3e-button
          variant="filled"
          .disabled="saving"
          @click="onPrimaryAction"
        >
          {{
            !isEdit && !isLast
              ? t('pages.profiles.wizard.next')
              : t('pages.profiles.wizard.save')
          }}
        </m3e-button>
      </div>
    </div>
  </m3e-dialog>
</template>

<style>
.character-wizard-dialog {
  --m3e-dialog-max-width: min(36rem, 92vw);
}
</style>
