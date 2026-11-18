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
  <m3e-dialog
    class="timeline-clip-voice-dialog"
    .open="show"
    @closed="setShow(false)"
  >
    <span slot="header">{{ title }}</span>
    <VoiceEditor
      :model-value="draft.voice"
      @update:model-value="onVoiceUpdate"
    />
    <div slot="actions" class="flex w-full justify-end gap-2">
      <m3e-button variant="text" @click="setShow(false)">
        {{ t('pages.generate.timeline.dialogCancel') }}
      </m3e-button>
      <m3e-button variant="filled" @click="onConfirm">
        {{ t('pages.generate.timeline.dialogApply') }}
      </m3e-button>
    </div>
  </m3e-dialog>
</template>

<style>
.timeline-clip-voice-dialog {
  --m3e-dialog-max-width: min(26rem, 92vw);
}
</style>
