<script setup lang="ts">
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import type { TimelineClipSpeaker } from '../../types/timeline';
import type { Voice } from '../../backend/types';
import { cloneSpeaker } from '../../composables/timeline/clipModel';
import VoiceEditor from '../voice/VoiceEditor.vue';

const props = defineProps<{
  show: boolean;
  speaker: TimelineClipSpeaker;
  title: string;
}>();

const emit = defineEmits<{
  'update:show': [value: boolean];
  confirm: [speaker: TimelineClipSpeaker];
}>();

const { t } = useI18n();

const draft = ref<TimelineClipSpeaker>(cloneSpeaker(props.speaker));

watch(
  () => props.show,
  (show) => {
    if (show) draft.value = cloneSpeaker(props.speaker);
  },
);

function setShow(value: boolean): void {
  emit('update:show', value);
}

function onVoiceUpdate(voice: Voice): void {
  draft.value = { ...draft.value, voice };
}

function onConfirm(): void {
  emit('confirm', cloneSpeaker(draft.value));
  setShow(false);
}
</script>

<template>
  <var-dialog
    :show="show"
    dialog-class="timeline-clip-voice-dialog"
    width="26rem"
    :title="title"
    :confirm-button="false"
    :cancel-button="false"
    :close-on-click-overlay="true"
    :close-on-key-escape="true"
    @update:show="setShow"
  >
    <VoiceEditor
      :model-value="draft.voice"
      @update:model-value="onVoiceUpdate"
    />
    <template #actions="{ slotClass }">
      <div :class="[slotClass, 'flex w-full justify-end gap-2']">
        <var-button
          text
          @click="setShow(false)"
        >
          {{ t('pages.generate.timeline.dialogCancel') }}
        </var-button>
        <var-button
          type="primary"
          @click="onConfirm"
        >
          {{ t('pages.generate.timeline.dialogApply') }}
        </var-button>
      </div>
    </template>
  </var-dialog>
</template>

<style>
.timeline-clip-voice-dialog {
  max-width: min(26rem, 92vw);
}
</style>
