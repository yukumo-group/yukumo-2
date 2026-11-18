<script setup lang="ts">
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { Play, RotateCcw, Save } from '@lucide/vue';
import { Snackbar } from '@varlet/ui';
import SettingsCard from '../settings/SettingsCard.vue';
import VersionSelectLabel from '../profiles/VersionSelectLabel.vue';
import VoiceConfigPanel from '../voice/VoiceConfigPanel.vue';
import SpeechTextEditor from '../voice/SpeechTextEditor.vue';
import { createDefaultAquesTalk10 } from '../../composables/aquestalkPresets';
import {
  useYukumoJsTalk,
  type YukumoJsTalkInput,
} from '../../composables/useYukumoJsTalk';
import type {
  AquesTalk1Config,
  AquesTalk10Config,
  AquesTalk2Config,
  AquesTalkVersion,
} from '../../types/profile';
import {
  CUSTOM_PROFILE_ID,
  createDefaultAquesTalk1,
  createDefaultAquesTalk2,
} from '../../types/profile';

const { t } = useI18n();
const { busy, preview, save, stopPlayback } = useYukumoJsTalk();

const version = ref<AquesTalkVersion>(2);
const profileId = ref(CUSTOM_PROFILE_ID);
const text = ref('');
const at1 = ref<AquesTalk1Config>(createDefaultAquesTalk1());
const at2 = ref<AquesTalk2Config>(createDefaultAquesTalk2());
const at10 = ref<AquesTalk10Config>(createDefaultAquesTalk10());

const versionLabelKey: Record<AquesTalkVersion, string> = {
  1: 'pages.utilities.quickGenerate.versions.at1',
  2: 'pages.utilities.quickGenerate.versions.at2',
  10: 'pages.utilities.quickGenerate.versions.at10',
};

const selectedVersionLabel = computed(() => t(versionLabelKey[version.value]));

const canGenerate = computed(() => text.value.trim().length > 0);
const actionsDisabled = computed(
  () => !canGenerate.value || busy.value != null,
);

function talkInput(): YukumoJsTalkInput {
  return {
    version: version.value,
    text: text.value,
    at1: at1.value,
    at2: at2.value,
    at10: at10.value,
  };
}

function reportError(error: unknown): void {
  Snackbar({
    content: t('pages.utilities.quickGenerate.failed', {
      error: error instanceof Error ? error.message : String(error),
    }),
    type: 'error',
  });
}

function onVersionChange(value: string | number): void {
  const n = typeof value === 'string' ? Number(value) : value;
  if (n !== 1 && n !== 2 && n !== 10) return;
  version.value = n;
}

async function onPreview(): Promise<void> {
  if (actionsDisabled.value) return;
  try {
    await preview(talkInput());
  } catch (error) {
    reportError(error);
  }
}

async function onSave(): Promise<void> {
  if (actionsDisabled.value) return;
  try {
    await save(talkInput());
  } catch (error) {
    reportError(error);
  }
}

function onReset(): void {
  stopPlayback();
  text.value = '';
  if (profileId.value === CUSTOM_PROFILE_ID) {
    at1.value = createDefaultAquesTalk1();
    at2.value = createDefaultAquesTalk2();
    at10.value = createDefaultAquesTalk10();
  }
}
</script>

<template>
  <SettingsCard>
    <div class="flex flex-col gap-4 p-4 sm:p-5">
      <header class="flex flex-col gap-1">
        <h2 class="font-bold text-text text-xl">
          {{ t('pages.utilities.quickGenerate.title') }}
        </h2>
        <p class="text-sm text-on-surface-variant">
          {{ t('pages.utilities.quickGenerate.description') }}
        </p>
      </header>

      <div
        class="flex min-h-0 flex-col gap-4 sm:flex-row sm:items-stretch"
      >
        <aside
          class="flex w-full shrink-0 flex-col gap-3 sm:w-72 sm:border-r sm:border-outline/30 sm:pr-4"
          :aria-label="t('pages.utilities.quickGenerate.configAriaLabel')"
        >
          <var-select
            :model-value="version"
            size="small"
            variant="outlined"
            :placeholder="t('pages.utilities.quickGenerate.version')"
            @update:model-value="onVersionChange"
          >
            <template #selected>
              <VersionSelectLabel
                :version="version"
                :label="selectedVersionLabel"
              />
            </template>
            <var-option
              :value="1"
              :label="t('pages.utilities.quickGenerate.versions.at1')"
            >
              <VersionSelectLabel
                :version="1"
                :label="t('pages.utilities.quickGenerate.versions.at1')"
              />
            </var-option>
            <var-option
              :value="2"
              :label="t('pages.utilities.quickGenerate.versions.at2')"
            >
              <VersionSelectLabel
                :version="2"
                :label="t('pages.utilities.quickGenerate.versions.at2')"
              />
            </var-option>
            <var-option
              :value="10"
              :label="t('pages.utilities.quickGenerate.versions.at10')"
            >
              <VersionSelectLabel
                :version="10"
                :label="t('pages.utilities.quickGenerate.versions.at10')"
              />
            </var-option>
          </var-select>

          <VoiceConfigPanel
            :version="version"
            v-model:profile-id="profileId"
            v-model:at1="at1"
            v-model:at2="at2"
            v-model:at10="at10"
          />
        </aside>

        <div class="flex min-h-0 min-w-0 flex-1 flex-col gap-3">
          <SpeechTextEditor
            v-model="text"
            :ariaLabel="t('pages.utilities.quickGenerate.editorAriaLabel')"
          />

          <div class="flex flex-wrap items-center gap-2">
            <var-button
              type="primary"
              :loading="busy === 'preview'"
              :disabled="actionsDisabled"
              @click="onPreview"
            >
              <Play
                :size="18"
                aria-hidden="true"
                class="mr-2"
              />
              {{ t('pages.utilities.quickGenerate.preview') }}
            </var-button>
            <var-button
              type="primary"
              outlined
              :loading="busy === 'save'"
              :disabled="actionsDisabled"
              @click="onSave"
            >
              <Save
                :size="18"
                aria-hidden="true"
                class="mr-2"
              />
              {{ t('pages.utilities.quickGenerate.save') }}
            </var-button>
            <div class="grow" />
            <var-button
              tonal
              @click="onReset"
            >
              <RotateCcw
                :size="18"
                aria-hidden="true"
                class="mr-2"
              />
              {{ t('pages.utilities.quickGenerate.reset') }}
            </var-button>
          </div>
        </div>
      </div>
    </div>
  </SettingsCard>
</template>
