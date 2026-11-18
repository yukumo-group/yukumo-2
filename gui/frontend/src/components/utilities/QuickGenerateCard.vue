<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, useId, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { Play, RotateCcw, Save } from '@lucide/vue';
import type { M3eSelectElement } from '@m3e/web/select';
import { notify } from '../../ui/notify';
import SettingsCard from '../settings/SettingsCard.vue';
import VersionSelectLabel from '../profiles/VersionSelectLabel.vue';
import SynthParamForm from '../voice/SynthParamForm.vue';
import SpeechTextEditor from '../voice/SpeechTextEditor.vue';
import { getBackend } from '../../backend';
import {
  createDefaultParam,
  createDefaultVoice,
  normalizeVoice,
  paramForVersion,
  upsertParam,
  voiceSupportsVersion,
  voiceVersions,
} from '../../backend/params';
import type { AquesTalkVersion, Voice } from '../../backend/types';
import { isAquesTalkVersion } from '../../backend/types';
import { downloadWav, wavFilename } from '../../backend/wav';
import { useCharacters } from '../../store/useCharacters';

const CUSTOM_SOURCE = '__custom__';
const DEFAULT_VERSIONS = [1, 2, 10] as const;

const { t } = useI18n();
const versionSelectId = useId();
const sourceSelectId = useId();
const { characters, ensureLoaded } = useCharacters();

const version = ref<AquesTalkVersion>(2);
const source = ref(CUSTOM_SOURCE);
const voice = ref<Voice>(createDefaultVoice('custom', DEFAULT_VERSIONS));
const text = ref('');
const busy = ref<'preview' | 'save' | null>(null);


const versionLabelKey: Record<AquesTalkVersion, string> = {
  1: 'pages.utilities.quickGenerate.versions.at1',
  2: 'pages.utilities.quickGenerate.versions.at2',
  10: 'pages.utilities.quickGenerate.versions.at10',
};

const voiceOptions = computed(() =>
  characters.value.flatMap((character) =>
    character.voices.map((item) => ({
      value: `${character.name}:${item.name}`,
      label: `${character.name} / ${item.name}`,
      voice: item,
    })),
  ),
);

const selectedSourceLabel = computed(() => {
  if (source.value === CUSTOM_SOURCE) {
    return t('pages.utilities.quickGenerate.customVoice');
  }
  return (
    voiceOptions.value.find((option) => option.value === source.value)?.label ??
    t('pages.utilities.quickGenerate.customVoice')
  );
});

const activeParam = computed(() => paramForVersion(voice.value, version.value));
const canGenerate = computed(
  () => text.value.trim().length > 0 && activeParam.value !== undefined,
);
const actionsDisabled = computed(
  () => !canGenerate.value || busy.value !== null,
);

onMounted(() => {
  void ensureLoaded().catch(reportError);
});

onBeforeUnmount(() => {
  void getBackend().stopPlayback();
});

watch(version, (next) => {
  if (paramForVersion(voice.value, next)) return;
  voice.value = upsertParam(voice.value, createDefaultParam(next));
});

watch(characters, () => {
  if (
    source.value !== CUSTOM_SOURCE &&
    !voiceOptions.value.some((option) => option.value === source.value)
  ) {
    source.value = CUSTOM_SOURCE;
  }
});

function reportError(cause: unknown): void {
  notify(
    t('pages.utilities.quickGenerate.failed', {
      error: cause instanceof Error ? cause.message : String(cause),
    }),
  );
}

function selectString(event: Event): string | null {
  const select = event.currentTarget as M3eSelectElement;
  return typeof select.value === 'string' ? select.value : null;
}

function cloneVoice(value: Voice): Voice {
  return {
    name: value.name,
    params: value.params.map((param) => ({ ...param })),
  };
}

function onVersionChange(event: Event): void {
  const value = selectString(event);
  if (value === null) return;
  const next = Number(value);
  if (!isAquesTalkVersion(next)) return;
  version.value = next;
}

function onSourceChange(event: Event): void {
  const next = selectString(event);
  if (next === null) return;
  source.value = next;
  if (next === CUSTOM_SOURCE) {
    voice.value = createDefaultVoice('custom', DEFAULT_VERSIONS);
    return;
  }
  const option = voiceOptions.value.find((item) => item.value === next);
  if (!option) return;
  const nextVoice = cloneVoice(option.voice);
  if (!voiceSupportsVersion(nextVoice, version.value)) {
    const available = voiceVersions(nextVoice);
    version.value = available[0] ?? version.value;
  }
  voice.value = nextVoice;
}

function onParamChange(value: Voice['params'][number]): void {
  voice.value = upsertParam(voice.value, value);
}

async function synthesize(): Promise<Uint8Array> {
  const param = activeParam.value;
  if (!param) throw new Error(t('pages.utilities.quickGenerate.noParam'));
  return getBackend().synthesize({
    voice: normalizeVoice(voice.value),
    text: text.value.trim(),
    version: version.value,
  });
}

async function onPreview(): Promise<void> {
  if (actionsDisabled.value) return;
  busy.value = 'preview';
  try {
    const wav = await synthesize();
    await getBackend().playWav(wav);
  } catch (cause) {
    reportError(cause);
  } finally {
    busy.value = null;
  }
}

async function onSave(): Promise<void> {
  if (actionsDisabled.value) return;
  busy.value = 'save';
  try {
    const wav = await synthesize();
    downloadWav(wav, wavFilename(text.value));
  } catch (cause) {
    reportError(cause);
  } finally {
    busy.value = null;
  }
}

function onReset(): void {
  void getBackend().stopPlayback();
  text.value = '';
  version.value = 2;
  source.value = CUSTOM_SOURCE;
  voice.value = createDefaultVoice('custom', DEFAULT_VERSIONS);
}
</script>

<template>
  <SettingsCard>
    <div class="flex flex-col gap-4 p-4 sm:p-5">
      <header class="flex flex-col gap-1">
        <h2 class="font-bold text-on-surface text-xl">
          {{ t('pages.utilities.quickGenerate.title') }}
        </h2>
        <p class="text-sm text-on-surface-variant">
          {{ t('pages.utilities.quickGenerate.description') }}
        </p>
        <p class="text-xs text-on-surface-variant">
          {{ t('pages.utilities.quickGenerate.inputHint') }}
        </p>
      </header>

      <div
        class="flex min-h-0 flex-col gap-4 sm:flex-row sm:items-stretch"
      >
        <aside
          class="flex w-full shrink-0 flex-col gap-3 sm:w-80 sm:border-r sm:border-outline/30 sm:pr-4"
          :aria-label="t('pages.utilities.quickGenerate.configAriaLabel')"
        >
          <m3e-form-field>
            <label slot="label" :for="versionSelectId">
              {{ t('pages.utilities.quickGenerate.version') }}
            </label>
            <m3e-select
              :id="versionSelectId"
              .value="String(version)"
              @change="onVersionChange"
            >
              <VersionSelectLabel
                slot="value"
                :version="version"
                :label="t(versionLabelKey[version])"
              />
              <m3e-option
                v-for="item in ([1, 2, 10] as const)"
                :key="item"
                :value="String(item)"
              >
                <VersionSelectLabel
                  :version="item"
                  :label="t(versionLabelKey[item])"
                />
              </m3e-option>
            </m3e-select>
          </m3e-form-field>

          <m3e-form-field>
            <label slot="label" :for="sourceSelectId">
              {{ t('pages.utilities.quickGenerate.voiceSource') }}
            </label>
            <m3e-select
              :id="sourceSelectId"
              .value="source"
              @change="onSourceChange"
            >
              <span slot="value" class="truncate">{{ selectedSourceLabel }}</span>
              <m3e-option :value="CUSTOM_SOURCE">
                {{ t('pages.utilities.quickGenerate.customVoice') }}
              </m3e-option>
              <m3e-option
                v-for="option in voiceOptions"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </m3e-option>
            </m3e-select>
          </m3e-form-field>

          <div class="flex max-h-96 flex-col gap-3 overflow-y-auto sm:max-h-none">
            <SynthParamForm
              v-if="activeParam"
              :model-value="activeParam"
              @update:model-value="onParamChange"
            />
            <p
              v-else
              class="text-sm text-on-surface-variant"
            >
              {{ t('pages.utilities.quickGenerate.noParam') }}
            </p>
          </div>
        </aside>

        <div class="flex min-h-0 min-w-0 flex-1 flex-col gap-3">
          <SpeechTextEditor
            v-model="text"
            :ariaLabel="t('pages.utilities.quickGenerate.editorAriaLabel')"
          />

          <div class="flex flex-wrap items-center gap-2">
            <m3e-button
              variant="filled"
              .disabled="actionsDisabled || busy === 'preview'"
              @click="onPreview"
            >
              <Play slot="icon" :size="18" aria-hidden="true" />
              {{ t('pages.utilities.quickGenerate.preview') }}
            </m3e-button>
            <m3e-button
              variant="outlined"
              .disabled="actionsDisabled || busy === 'save'"
              @click="onSave"
            >
              <Save slot="icon" :size="18" aria-hidden="true" />
              {{ t('pages.utilities.quickGenerate.save') }}
            </m3e-button>
            <div class="grow" />
            <m3e-button
              variant="tonal"
              @click="onReset"
            >
              <RotateCcw slot="icon" :size="18" aria-hidden="true" />
              {{ t('pages.utilities.quickGenerate.reset') }}
            </m3e-button>
          </div>
        </div>
      </div>
    </div>
  </SettingsCard>
</template>
