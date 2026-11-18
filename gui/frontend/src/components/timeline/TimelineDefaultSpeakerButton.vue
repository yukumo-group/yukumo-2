<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import type { TimelineClipSpeaker } from '../../types/timeline';
import { User } from '@lucide/vue';
import Hint from '../ui/Hint.vue';

defineProps<{
  speaker: TimelineClipSpeaker;
}>();

const emit = defineEmits<{
  click: [];
}>();

const { t } = useI18n();
</script>

<template>
  <Hint
    :text="t('pages.generate.timeline.defaultSpeaker')"
    position="below"
    v-slot="{ id }"
  >
    <button
      :id="id"
      type="button"
      class="inline-flex size-8 shrink-0 items-center justify-center rounded-full hover:bg-surface-container-high"
      :aria-label="t('pages.generate.timeline.defaultSpeaker')"
      @click="emit('click')"
    >
      <span
        class="flex size-6 items-center justify-center overflow-hidden rounded-full bg-surface-container-high text-on-surface-variant"
      >
        <User
          :size="14"
          aria-hidden="true"
        />
      </span>
    </button>
  </Hint>
</template>
