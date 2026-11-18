<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import type { AquesTalkVersion } from '../../../types/profile';

const name = defineModel<string>('name', { required: true });
const description = defineModel<string>('description', { required: true });
const version = defineModel<AquesTalkVersion>('version', { required: true });

const emit = defineEmits<{
  clearError: [];
}>();

const { t } = useI18n();

const versionOptions: {
  value: AquesTalkVersion;
  labelKey: string;
  dotClass: string;
}[] = [
  {
    value: 1,
    labelKey: 'pages.profiles.engines.at1',
    dotClass: 'bg-success',
  },
  {
    value: 2,
    labelKey: 'pages.profiles.engines.at2',
    dotClass: 'bg-warning',
  },
  {
    value: 10,
    labelKey: 'pages.profiles.engines.at10',
    dotClass: 'bg-danger',
  },
];
</script>

<template>
  <var-input
    v-model="name"
    size="small"
    variant="outlined"
    :placeholder="t('pages.profiles.wizard.namePlaceholder')"
    :aria-label="t('pages.profiles.wizard.namePlaceholder')"
    @update:model-value="emit('clearError')"
  />
  <var-input
    v-model="description"
    size="small"
    variant="outlined"
    textarea
    :rows="3"
    :placeholder="t('pages.profiles.wizard.descriptionPlaceholder')"
    :aria-label="t('pages.profiles.wizard.descriptionPlaceholder')"
  />
  <div class="flex items-center justify-between gap-3">
    <p class="shrink-0 text-base font-medium text-text">
      {{ t('pages.profiles.wizard.version') }}
    </p>
    <var-segmented-buttons
      v-model="version"
      checkmark
      size="small"
      @update:model-value="emit('clearError')"
    >
      <var-segmented-button
        v-for="option in versionOptions"
        :key="option.value"
        :checked-value="option.value"
      >
        <template #default="{ checked }">
          <span class="inline-flex items-center gap-1.5">
            <span
              v-if="!checked"
              class="size-1.5 shrink-0 rounded-full"
              :class="option.dotClass"
              aria-hidden="true"
            />
            {{ t(option.labelKey) }}
          </span>
        </template>
      </var-segmented-button>
    </var-segmented-buttons>
  </div>
</template>
