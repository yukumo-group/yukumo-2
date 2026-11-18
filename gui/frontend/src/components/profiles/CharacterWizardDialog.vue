<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { Dialog } from '@varlet/ui';
import { useI18n } from 'vue-i18n';
import { Plus } from '@lucide/vue';
import VoiceEditor from '../voice/VoiceEditor.vue';
import {
  createDefaultVoice,
  normalizeVoice,
} from '../../backend/params';
import type { Character, Voice } from '../../backend/types';
import { useCharacters } from '../../store/useCharacters';

type WizardStepKind = 'basics' | 'voices';

interface SwipeExpose {
  to: (index: number) => void;
}

const props = defineProps<{
  show: boolean;
  character: Character | null;
}>();

const emit = defineEmits<{
  'update:show': [value: boolean];
}>();

const { t } = useI18n();
const { addCharacter, updateCharacter, deleteCharacter } = useCharacters();

const swipeRef = ref<SwipeExpose | null>(null);
const swipeSession = ref(0);
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
  swipeSession.value += 1;
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
  const next = Math.max(0, Math.min(index, stepKinds.value.length - 1));
  stepIndex.value = next;
  swipeRef.value?.to(next);
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

  const result = await Dialog({
    title: t('pages.profiles.wizard.deleteConfirmTitle'),
    message: t('pages.profiles.wizard.deleteConfirmMessage', {
      name: character.name,
    }),
    confirmButtonText: t('pages.profiles.wizard.deleteConfirm'),
    cancelButtonText: t('pages.profiles.wizard.deleteCancel'),
    confirmButtonProps: { type: 'danger' },
  });

  if (result !== 'confirm') return;
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
  <var-dialog
    :show="show"
    dialog-class="character-wizard-dialog"
    width="36rem"
    :title="
      isEdit
        ? t('pages.profiles.wizard.editTitle')
        : t('pages.profiles.wizard.createTitle')
    "
    :confirm-button="false"
    :cancel-button="false"
    :close-on-click-overlay="true"
    :close-on-key-escape="true"
    @update:show="setShow"
  >
    <div class="flex flex-col gap-3 text-left">
      <var-steps
        :active="stepIndex"
        @click-step="onClickStep"
      >
        <var-step
          v-for="kind in stepKinds"
          :key="kind"
        >
          {{ t(stepLabelKey[kind]) }}
        </var-step>
      </var-steps>

      <p
        v-if="error"
        class="rounded-lg bg-danger/15 px-3 py-2 text-base font-medium text-danger"
        role="alert"
      >
        {{ error }}
      </p>

      <var-swipe
        :key="swipeSession"
        ref="swipeRef"
        class="wizard-swipe"
        :loop="false"
        :touchable="false"
        :indicator="false"
        :initial-index="stepIndex"
      >
        <var-swipe-item
          v-for="kind in stepKinds"
          :key="kind"
        >
          <div class="wizard-pane">
            <template v-if="kind === 'basics'">
              <var-input
                v-model="name"
                size="small"
                variant="outlined"
                :placeholder="t('pages.profiles.wizard.namePlaceholder')"
                :aria-label="t('pages.profiles.wizard.namePlaceholder')"
                @update:model-value="clearError"
              />
              <var-input
                v-model="description"
                size="small"
                variant="outlined"
                textarea
                :rows="3"
                :placeholder="t('pages.profiles.wizard.descriptionPlaceholder')"
                :aria-label="t('pages.profiles.wizard.descriptionPlaceholder')"
              />
            </template>

            <template v-else>
              <div class="flex items-center justify-between gap-2">
                <p class="text-base font-medium text-text">
                  {{ t('pages.profiles.wizard.voices') }}
                </p>
                <var-button
                  size="small"
                  outlined
                  @click="addVoice"
                >
                  <Plus
                    :size="14"
                    aria-hidden="true"
                    class="mr-1"
                  />
                  {{ t('pages.profiles.wizard.addVoice') }}
                </var-button>
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
            </template>
          </div>
        </var-swipe-item>
      </var-swipe>
    </div>

    <template #actions="{ slotClass }">
      <div :class="[slotClass, 'flex w-full items-center gap-2']">
        <var-button
          v-if="isEdit"
          type="danger"
          text
          @click="onDelete"
        >
          {{ t('pages.profiles.wizard.delete') }}
        </var-button>
        <div class="ml-auto flex gap-2">
          <var-button
            text
            @click="cancel"
          >
            {{ t('pages.profiles.wizard.cancel') }}
          </var-button>
          <var-button
            v-if="!isEdit"
            text
            :disabled="stepIndex === 0"
            @click="goBack"
          >
            {{ t('pages.profiles.wizard.back') }}
          </var-button>
          <var-button
            type="primary"
            :loading="saving"
            @click="onPrimaryAction"
          >
            {{
              !isEdit && !isLast
                ? t('pages.profiles.wizard.next')
                : t('pages.profiles.wizard.save')
            }}
          </var-button>
        </div>
      </div>
    </template>
  </var-dialog>
</template>

<style scoped>
.wizard-swipe {
  height: 26rem;
}

.wizard-pane {
  display: flex;
  height: 100%;
  flex-direction: column;
  gap: 0.75rem;
  overflow-y: auto;
  padding: 0.5rem 1rem 0.75rem;
}
</style>

<style>
.character-wizard-dialog {
  max-width: min(36rem, 92vw);
  --dialog-actions-padding: 0 24px 24px;
}

.character-wizard-dialog .var-dialog__message {
  max-height: min(36rem, 75vh);
  overflow: auto;
}

.character-wizard-dialog .var-steps {
  overflow-x: auto;
}
</style>
