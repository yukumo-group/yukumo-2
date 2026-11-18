<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { RouterView } from 'vue-router';
import { Terminal, X } from '@lucide/vue';
import Hint from '../components/ui/Hint.vue';
import AppRailNav from '../components/navigation/AppRailNav.vue';
import AppBottomNav from '../components/navigation/AppBottomNav.vue';
import { useM3eSplit } from '../composables/useM3eSplit';
import {
  TERMINAL_HEIGHT_MAX,
  TERMINAL_HEIGHT_MIN,
  TERMINAL_MAIN_RESERVE,
  latestStatusMessage,
  showTerminal,
  terminalHeight,
  toggleTerminal,
} from '../composables/terminal';

const { t } = useI18n();

const isSmUp = ref(false);

const tooltipPosition = computed(() => (isSmUp.value ? 'below' : 'above'));

const { splitRef, onSplitInput, onSplitChange } = useM3eSplit({
  open: showTerminal,
  size: terminalHeight,
  side: 'end',
  orientation: 'vertical',
  minPx: TERMINAL_HEIGHT_MIN,
  maxPx: TERMINAL_HEIGHT_MAX,
  reservePx: TERMINAL_MAIN_RESERVE,
});

const statusMessage = computed(
  () => latestStatusMessage.value ?? t('layout.statusIdle'),
);

const terminalToggleLabel = computed(() =>
  showTerminal.value
    ? t('layout.closeTerminal')
    : t('layout.toggleTerminal'),
);

let smMediaQuery: MediaQueryList | undefined;

function syncViewport(): void {
  if (!smMediaQuery) return;
  isSmUp.value = smMediaQuery.matches;
}

onMounted(() => {
  smMediaQuery = window.matchMedia('(min-width: 640px)');
  syncViewport();
  smMediaQuery.addEventListener('change', syncViewport);
});

onUnmounted(() => {
  smMediaQuery?.removeEventListener('change', syncViewport);
});
</script>

<template>
  <div
    class="layout-chrome relative flex h-screen w-screen flex-col overflow-hidden sm:flex-row"
  >
    <AppRailNav />

    <main
      class="relative z-10 flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden p-0 sm:p-1"
    >
      <div
        class="flex min-h-0 min-w-0 flex-1 flex-col items-center overflow-hidden rounded-none bg-surface p-0 sm:rounded-2xl sm:p-2"
      >
        <div class="flex min-h-0 w-full min-w-0 flex-1 flex-col">
          <m3e-split-pane
            ref="splitRef"
            class="min-h-0 w-full flex-1"
            :class="{ 'split-collapsed': !showTerminal }"
            orientation="vertical"
            .disabled="!showTerminal"
            :label="t('layout.resizeTerminal')"
            @input="onSplitInput"
            @change="onSplitChange"
          >
            <div
              slot="start"
              class="flex h-full min-h-0 w-full min-w-0 flex-col items-center overflow-x-hidden overflow-y-auto"
            >
              <RouterView v-slot="{ Component, route }">
                <Transition name="page">
                  <component :is="Component" :key="route.path" />
                </Transition>
              </RouterView>
            </div>
            <aside
              slot="end"
              class="h-full min-h-0 overflow-y-auto rounded-xl bg-surface-container p-3 text-on-surface-variant text-sm"
              :aria-hidden="!showTerminal"
              :aria-label="t('layout.terminalAriaLabel')"
            >
              {{ t('layout.terminalEmpty') }}
            </aside>
          </m3e-split-pane>

          <div
            class="flex shrink-0 items-center gap-1 border-outline/30 px-2 py-0.5"
            :class="showTerminal ? 'border-t' : ''"
            role="status"
            :aria-label="t('layout.statusBarAriaLabel')"
          >
            <p class="min-w-0 flex-1 truncate px-1 text-on-surface-variant text-sm">
              {{ statusMessage }}
            </p>

            <Hint
              :text="terminalToggleLabel"
              :position="tooltipPosition"
            >
              <template #default="{ id }">
                <m3e-icon-button
                  :id="id"
                  :aria-pressed="showTerminal"
                  :aria-label="terminalToggleLabel"
                  @click="toggleTerminal"
                >
                  <X
                    v-if="showTerminal"
                    :size="18"
                    aria-hidden="true"
                  />
                  <Terminal
                    v-else
                    :size="18"
                    aria-hidden="true"
                  />
                </m3e-icon-button>
              </template>
            </Hint>
          </div>
        </div>
      </div>
    </main>

    <AppBottomNav />
  </div>
</template>

<style scoped>
.layout-chrome {
  background-color: var(--md-sys-color-surface-container-low);
}
</style>
