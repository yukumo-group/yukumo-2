<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { Dialog } from '@varlet/ui';
import { useI18n } from 'vue-i18n';
import { createDefaultAquesTalk10 } from '../../composables/aquestalkPresets';
import {
  createProfile,
  removeProfile,
  updateProfile,
} from '../../composables/useCharacterProfiles';
import type {
  AquesTalk1Config,
  AquesTalk10Config,
  AquesTalk2Config,
  AquesTalkVersion,
  CharacterProfile,
} from '../../types/profile';
import {
  createDefaultAquesTalk1,
  createDefaultAquesTalk2,
} from '../../types/profile';
import WizardAt1Step from './wizard/WizardAt1Step.vue';
import WizardAt2Step from './wizard/WizardAt2Step.vue';
import WizardAt10Step from './wizard/WizardAt10Step.vue';
import WizardBasicsStep from './wizard/WizardBasicsStep.vue';
import WizardImageStep from './wizard/WizardImageStep.vue';

type WizardStepKind = 'basics' | 'voice' | 'image';

interface SwipeExpose {
  to: (index: number) => void;
}

const props = defineProps<{
  show: boolean;
  profile: CharacterProfile | null;
}>();

const emit = defineEmits<{
  'update:show': [value: boolean];
}>();

const { t } = useI18n();

const swipeRef = ref<SwipeExpose | null>(null);
const swipeSession = ref(0);
const stepIndex = ref(0);
const error = ref('');
const name = ref('');
const description = ref('');
const imageDataUrl = ref<string | null>(null);
const version = ref<AquesTalkVersion>(2);
const at1 = ref<AquesTalk1Config>(createDefaultAquesTalk1());
const at2 = ref<AquesTalk2Config>(createDefaultAquesTalk2());
const at10 = ref<AquesTalk10Config>(createDefaultAquesTalk10());

const isEdit = computed(() => props.profile !== null);

const stepKinds = computed<WizardStepKind[]>(() => ['basics', 'voice', 'image']);

const currentKind = computed(
  () => stepKinds.value[stepIndex.value] ?? 'basics',
);

const isLast = computed(
  () => stepIndex.value >= stepKinds.value.length - 1,
);

const stepLabelKey: Record<WizardStepKind, string> = {
  basics: 'pages.profiles.wizard.steps.basics',
  voice: 'pages.profiles.wizard.steps.voice',
  image: 'pages.profiles.wizard.steps.image',
};

watch(
  () => props.show,
  (open) => {
    if (!open) return;
    hydrate();
    stepIndex.value = 0;
    error.value = '';
    swipeSession.value += 1;
  },
);

function hydrate(): void {
  const profile = props.profile;
  if (!profile) {
    name.value = '';
    description.value = '';
    imageDataUrl.value = null;
    version.value = 2;
    Object.assign(at1.value, createDefaultAquesTalk1());
    Object.assign(at2.value, createDefaultAquesTalk2());
    Object.assign(at10.value, createDefaultAquesTalk10());
    return;
  }

  name.value = profile.name;
  description.value = profile.description;
  imageDataUrl.value = profile.imageDataUrl;
  version.value = profile.version;
  Object.assign(at1.value, profile.aquestalk1 ?? createDefaultAquesTalk1());
  Object.assign(at2.value, profile.aquestalk2 ?? createDefaultAquesTalk2());
  Object.assign(at10.value, profile.aquestalk10 ?? createDefaultAquesTalk10());
}

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
  }
  if (kind === 'voice') {
    if (version.value === 1 && !at1.value.voiceId.trim()) {
      error.value = t('pages.profiles.wizard.voiceIdRequired');
      return false;
    }
    if (version.value === 2 && !at2.value.phontName.trim()) {
      error.value = t('pages.profiles.wizard.phontNameRequired');
      return false;
    }
  }
  return true;
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
  onSave();
}

function onSave(): void {
  for (const kind of stepKinds.value) {
    if (!validateKind(kind)) {
      goTo(stepKinds.value.indexOf(kind));
      return;
    }
  }
  save();
}

function save(): void {
  const selectedVersion = version.value;
  const draft = {
    name: name.value.trim(),
    description: description.value.trim(),
    imageDataUrl: imageDataUrl.value,
    version: selectedVersion,
    aquestalk1:
      selectedVersion === 1
        ? { ...at1.value, voiceId: at1.value.voiceId.trim() }
        : undefined,
    aquestalk2:
      selectedVersion === 2
        ? { ...at2.value, phontName: at2.value.phontName.trim() }
        : undefined,
    aquestalk10: selectedVersion === 10 ? { ...at10.value } : undefined,
  };

  if (props.profile) {
    updateProfile(props.profile.id, draft);
  } else {
    createProfile(draft);
  }
  setShow(false);
}

async function onDelete(): Promise<void> {
  const profile = props.profile;
  if (!profile) return;

  const result = await Dialog({
    title: t('pages.profiles.wizard.deleteConfirmTitle'),
    message: t('pages.profiles.wizard.deleteConfirmMessage', {
      name: profile.name,
    }),
    confirmButtonText: t('pages.profiles.wizard.deleteConfirm'),
    cancelButtonText: t('pages.profiles.wizard.deleteCancel'),
    confirmButtonProps: { type: 'danger' },
  });

  if (result !== 'confirm') return;
  removeProfile(profile.id);
  setShow(false);
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
            <WizardBasicsStep
              v-if="kind === 'basics'"
              v-model:name="name"
              v-model:description="description"
              v-model:version="version"
              @clear-error="clearError"
            />
            <WizardAt1Step
              v-else-if="kind === 'voice' && version === 1"
              v-model="at1"
              @clear-error="clearError"
            />
            <WizardAt2Step
              v-else-if="kind === 'voice' && version === 2"
              v-model="at2"
              @clear-error="clearError"
            />
            <WizardAt10Step
              v-else-if="kind === 'voice' && version === 10"
              v-model="at10"
            />
            <WizardImageStep
              v-else-if="kind === 'image'"
              v-model="imageDataUrl"
              :active="currentKind === 'image'"
            />
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
  height: 14rem;
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
  max-height: min(32rem, 70vh);
  overflow: auto;
}

.character-wizard-dialog .var-steps {
  overflow-x: auto;
}
</style>
