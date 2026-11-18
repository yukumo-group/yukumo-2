<script setup lang="ts">
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import SpeechTextEditor from '../voice/SpeechTextEditor.vue';

const props = defineProps<{
  show: boolean;
  text: string;
}>();

const emit = defineEmits<{
  'update:show': [value: boolean];
  confirm: [text: string];
}>();

const { t } = useI18n();

const draft = ref(props.text);

watch(
  () => props.show,
  (show) => {
    if (show) draft.value = props.text;
  },
);

function setShow(value: boolean): void {
  emit('update:show', value);
}

function onConfirm(): void {
  emit('confirm', draft.value);
  setShow(false);
}
</script>

<template>
  <m3e-dialog
    class="timeline-clip-content-dialog"
    .open="show"
    @closed="setShow(false)"
  >
    <span slot="header">{{ t('pages.generate.timeline.contentDialogTitle') }}</span>
    <SpeechTextEditor
      v-model="draft"
      :ariaLabel="t('pages.utilities.quickGenerate.editorAriaLabel')"
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
.timeline-clip-content-dialog {
  --m3e-dialog-max-width: min(36rem, 92vw);
}
</style>
