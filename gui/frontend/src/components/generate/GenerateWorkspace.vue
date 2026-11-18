<script setup lang="ts">
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import type { editor as MonacoEditor } from 'monaco-editor';
import TimelineEditor from '../timeline/TimelineEditor.vue';
import { useM3eSplit } from '../../composables/useM3eSplit';
import {
  FILE_TREE_WIDTH_MAX,
  FILE_TREE_WIDTH_MIN,
  fileTreeWidth,
  rightPanel,
  showFileTree,
} from '../../composables/generateWorkspace';
import { resolvedTheme } from '../../theme/theme';

const { t } = useI18n();

const script = ref('');

const { splitRef, onSplitInput, onSplitChange } = useM3eSplit({
  open: showFileTree,
  size: fileTreeWidth,
  side: 'start',
  orientation: 'horizontal',
  minPx: FILE_TREE_WIDTH_MIN,
  maxPx: FILE_TREE_WIDTH_MAX,
  reservePx: 160,
});

const monacoTheme = computed(() =>
  resolvedTheme.value === 'dark' ? 'vs-dark' : 'vs',
);

const editorOptions: MonacoEditor.IStandaloneEditorConstructionOptions = {
  automaticLayout: true,
  fontSize: 14,
  minimap: { enabled: false },
  scrollBeyondLastLine: false,
  wordWrap: 'on',
  padding: { top: 12, bottom: 12 },
};
</script>

<template>
  <m3e-split-pane
    ref="splitRef"
    class="h-full min-h-0 min-w-0 flex-1"
    :class="{ 'split-collapsed': !showFileTree }"
    orientation="horizontal"
    .disabled="!showFileTree"
    :label="t('pages.generate.resizeFileTree')"
    @input="onSplitInput"
    @change="onSplitChange"
  >
    <aside
      slot="start"
      class="h-full min-h-0 overflow-hidden rounded-xl bg-surface-container"
      :aria-hidden="!showFileTree"
      :aria-label="t('pages.generate.fileTreeAriaLabel')"
    />

    <div slot="end" class="flex h-full w-full min-h-0 min-w-0 overflow-hidden">
    <div
      v-show="rightPanel === 'code'"
      class="min-h-0 min-w-0 flex-1 overflow-hidden rounded-xl border-2 border-outline/40 bg-surface-container"
      role="region"
      :aria-label="t('pages.generate.codeEditorAriaLabel')"
    >
      <VueMonacoEditor
        v-model:value="script"
        language="plaintext"
        :theme="monacoTheme"
        :options="editorOptions"
        height="100%"
      >
        <template #default>
          <span class="text-on-surface-variant text-sm">
            {{ t('pages.generate.editorLoading') }}
          </span>
        </template>
        <template #failure>
          <span class="text-error text-sm">
            {{ t('pages.generate.editorLoadFailed') }}
          </span>
        </template>
      </VueMonacoEditor>
    </div>

    <aside
      v-show="rightPanel === 'timeline'"
      class="min-h-0 min-w-0 flex-1 overflow-hidden rounded-xl bg-surface-container"
      role="region"
      :aria-label="t('pages.generate.timelineAriaLabel')"
    >
      <TimelineEditor />
    </aside>
    </div>
  </m3e-split-pane>
</template>
