<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue';
import { useId } from 'vue';
import { ChevronRight, Monitor, Moon, Sun } from '@lucide/vue';
import { useI18n } from 'vue-i18n';
import type { M3eSegmentedButtonElement } from '@m3e/web/segmented-button';
import type { M3eSelectElement } from '@m3e/web/select';
import type { M3eSliderThumbElement } from '@m3e/web/slider';
import { confirmDialog } from '../ui/confirm';
import { useRouter } from 'vue-router';
import SettingsCard from '../components/settings/SettingsCard.vue';
import SettingsItem from '../components/settings/SettingsItem.vue';
import {
  appLocale,
  isLocalePreference,
  resetLocale,
  setLocale,
  type LocalePreference,
} from '../i18n';
import { helpDestination } from '../navigation/destinations';
import {
  resetTheme,
  setHue,
  setPreference,
  themeHue,
  themePreference,
} from '../theme/theme';

const HUE_THROTTLE_MS = 120;

const { t } = useI18n();
const router = useRouter();
const localeSelectId = useId();

let lastPointer = { x: 0, y: 0 };
let pendingHue: number | null = null;
let hueThrottleTimer: ReturnType<typeof setTimeout> | null = null;

const localeOptions: { value: LocalePreference; labelKey: string }[] = [
  { value: 'system', labelKey: 'pages.settings.language.system' },
  { value: 'en-US', labelKey: 'pages.settings.language.enUS' },
  { value: 'ja-JP', labelKey: 'pages.settings.language.jaJP' },
  { value: 'zh-CN', labelKey: 'pages.settings.language.zhCN' },
];

const accentOpen = ref(false);
const sliderHue = ref(themeHue.value);

const accentPreview = computed(
  () => `hsl(${sliderHue.value} 75% 50%)`,
);

watch(themeHue, (value) => {
  if (pendingHue === null) {
    sliderHue.value = value;
  }
});

onBeforeUnmount(() => {
  clearHueThrottle();
  flushPendingHue();
});

function openHelp(): void {
  void router.push(helpDestination.path);
}

function capturePointer(event: Event): void {
  if (event instanceof MouseEvent || event instanceof PointerEvent) {
    lastPointer = { x: event.clientX, y: event.clientY };
    return;
  }
  if (typeof TouchEvent !== 'undefined' && event instanceof TouchEvent) {
    const touch = event.changedTouches[0] ?? event.touches[0];
    if (touch) {
      lastPointer = { x: touch.clientX, y: touch.clientY };
    }
  }
}

function onThemeChange(event: Event): void {
  const group = event.currentTarget as M3eSegmentedButtonElement;
  const value = group.value;
  if (value === 'light' || value === 'dark' || value === 'system') {
    void setPreference(value, lastPointer);
  }
}

function onHueInput(event: Event): void {
  const thumb = event.target as M3eSliderThumbElement;
  const next = thumb.value;
  if (typeof next !== 'number') {
    return;
  }

  sliderHue.value = next;
  pendingHue = next;

  if (hueThrottleTimer !== null) {
    return;
  }

  hueThrottleTimer = setTimeout(() => {
    hueThrottleTimer = null;
    flushPendingHue();
  }, HUE_THROTTLE_MS);
}

function onHueEnd(): void {
  clearHueThrottle();
  flushPendingHue();
}

function flushPendingHue(): void {
  if (pendingHue === null) {
    return;
  }
  setHue(pendingHue);
  pendingHue = null;
}

function clearHueThrottle(): void {
  if (hueThrottleTimer === null) {
    return;
  }
  clearTimeout(hueThrottleTimer);
  hueThrottleTimer = null;
}

function toggleAccentSlider(): void {
  accentOpen.value = !accentOpen.value;
}

function onLocaleChange(event: Event): void {
  const select = event.currentTarget as M3eSelectElement;
  const next = typeof select.value === 'string' ? select.value : '';
  if (isLocalePreference(next)) {
    setLocale(next);
  }
}

async function onResetClick(): Promise<void> {
  const confirmed = await confirmDialog({
    title: t('pages.settings.reset.confirmTitle'),
    message: t('pages.settings.reset.confirmMessage'),
    confirmText: t('pages.settings.reset.confirm'),
    cancelText: t('pages.settings.reset.cancel'),
    danger: true,
  });

  if (!confirmed) {
    return;
  }

  clearHueThrottle();
  pendingHue = null;
  accentOpen.value = false;
  resetTheme();
  resetLocale();
  sliderHue.value = themeHue.value;
}
</script>

<template>
  <section class="flex w-full max-w-3xl flex-col gap-6 px-6 pt-8 text-left sm:px-6">
    <h1 class="font-bold tracking-wide text-on-surface text-3xl">
      {{ t('pages.settings.title') }}
    </h1>

    <div class="flex flex-col gap-4">
      <SettingsCard>
        <SettingsItem
          :title="t('pages.settings.appearance.title')"
          :description="t('pages.settings.appearance.description')"
        >
        <m3e-segmented-button
          @click="capturePointer"
          @change="onThemeChange"
        >
          <m3e-button-segment
            value="light"
            .checked="themePreference === 'light'"
          >
            <Sun
              v-if="themePreference !== 'light'"
              slot="icon"
              :size="16"
              aria-hidden="true"
            />
            {{ t('pages.settings.appearance.light') }}
          </m3e-button-segment>
          <m3e-button-segment
            value="dark"
            .checked="themePreference === 'dark'"
          >
            <Moon
              v-if="themePreference !== 'dark'"
              slot="icon"
              :size="16"
              aria-hidden="true"
            />
            {{ t('pages.settings.appearance.dark') }}
          </m3e-button-segment>
          <m3e-button-segment
            value="system"
            .checked="themePreference === 'system'"
          >
            <Monitor
              v-if="themePreference !== 'system'"
              slot="icon"
              :size="16"
              aria-hidden="true"
            />
            {{ t('pages.settings.appearance.system') }}
          </m3e-button-segment>
        </m3e-segmented-button>
      </SettingsItem>

      <div class="settings-item flex flex-col">
        <div class="flex items-center justify-between gap-3 px-4 py-3">
          <div class="flex min-w-0 flex-1 flex-col gap-0.5">
            <p class="truncate font-semibold text-on-surface text-base">
              {{ t('pages.settings.accentColor.title') }}
            </p>
            <p class="truncate text-on-surface-variant text-sm">
              {{ t('pages.settings.accentColor.description') }}
            </p>
          </div>

          <button
            type="button"
            class="accent-swatch size-8 shrink-0 rounded-full border border-outline bg-primary"
            :class="{ 'accent-swatch--open': accentOpen }"
            :aria-label="t('pages.settings.accentColor.toggleSlider')"
            :aria-expanded="accentOpen"
            @click="toggleAccentSlider"
          />
        </div>

        <!-- Grid 0fr→1fr grows after the slider has a height. -->
        <div
          class="accent-panel"
          :data-open="accentOpen"
        >
          <div class="accent-panel__inner">
            <div class="px-4 pb-3">
              <m3e-slider
                class="hue-slider w-full"
                .min="0"
                .max="360"
                .step="1"
                :aria-label="t('pages.settings.accentColor.ariaLabel')"
                @input="onHueInput"
                @change="onHueEnd"
              >
                <m3e-slider-thumb .value="sliderHue" />
              </m3e-slider>
            </div>
          </div>
        </div>
      </div>

      <SettingsItem
        :title="t('pages.settings.language.title')"
        :description="t('pages.settings.language.description')"
      >
        <m3e-form-field class="w-52">
          <label slot="label" :for="localeSelectId">
            {{ t('pages.settings.language.title') }}
          </label>
          <m3e-select
            :id="localeSelectId"
            .value="appLocale"
            @change="onLocaleChange"
          >
            <m3e-option
              v-for="option in localeOptions"
              :key="option.value"
              :value="option.value"
            >
              {{ t(option.labelKey) }}
            </m3e-option>
          </m3e-select>
        </m3e-form-field>
      </SettingsItem>
    </SettingsCard>

    <SettingsCard>
      <SettingsItem
        :title="t('pages.settings.reset.title')"
        :description="t('pages.settings.reset.description')"
      >
        <m3e-button
          class="tone-error"
          variant="filled"
          @click="onResetClick"
        >
          {{ t('pages.settings.reset.action') }}
        </m3e-button>
      </SettingsItem>
    </SettingsCard>

    <div class="sm:hidden">
      <m3e-list class="overflow-hidden rounded-2xl">
        <m3e-list-item @click="openHelp">
          <component
            :is="helpDestination.icon"
            slot="leading"
          />
          <button type="button" class="w-full text-left" @click="openHelp">
            {{ t(helpDestination.labelKey) }}
          </button>
          <ChevronRight slot="trailing" :size="20" />
        </m3e-list-item>
      </m3e-list>
    </div>
    </div>
  </section>
</template>

<style scoped>
.accent-panel {
  display: grid;
  grid-template-rows: 0fr;
  transition: grid-template-rows 220ms cubic-bezier(0.2, 0, 0, 1);
}

.accent-panel[data-open='true'] {
  grid-template-rows: 1fr;
}

@media (prefers-reduced-motion: reduce) {
  .accent-panel {
    transition: none;
  }
}

.accent-panel__inner {
  overflow: hidden;
  min-height: 0;
}

.accent-swatch {
  cursor: pointer;
  transition: box-shadow 160ms ease;
}

.accent-swatch--open {
  box-shadow: 0 0 0 2px var(--md-sys-color-primary);
}

.hue-slider {
  --m3e-slider-active-track-color: transparent;
  --m3e-slider-inactive-track-color: transparent;
  --m3e-slider-thumb-color: v-bind(accentPreview);
  border-radius: 999px;
  background: linear-gradient(
    to right,
    hsl(0, 85%, 55%),
    hsl(60, 85%, 55%),
    hsl(120, 85%, 45%),
    hsl(180, 85%, 45%),
    hsl(240, 85%, 55%),
    hsl(300, 85%, 55%),
    hsl(360, 85%, 55%)
  );
}
</style>
