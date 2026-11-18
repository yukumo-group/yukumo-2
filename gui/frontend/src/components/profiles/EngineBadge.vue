<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { isAquesTalkVersion } from '../../backend/types';

const props = defineProps<{
  version: number;
}>();

const { t } = useI18n();

const safeVersion = computed(() => (isAquesTalkVersion(props.version) ? props.version : 1));

const tone = computed(() => {
  const map = {
    1: 'primary',
    2: 'tertiary',
    10: 'error',
  } as const;
  return map[safeVersion.value];
});

const labelKey = computed(() => {
  const map = {
    1: 'pages.profiles.engines.at1',
    2: 'pages.profiles.engines.at2',
    10: 'pages.profiles.engines.at10',
  } as const;
  return map[safeVersion.value];
});
</script>

<template>
  <m3e-badge
    class="engine-badge"
    :data-tone="tone"
    size="medium"
  >
    {{ t(labelKey) }}
  </m3e-badge>
</template>

<style scoped>
.engine-badge {
  box-shadow: 0 1px 3px color-mix(in srgb, var(--md-sys-color-on-surface) 20%, transparent);
}

.engine-badge[data-tone='primary'] {
  --m3e-badge-container-color: var(--md-sys-color-primary);
  --m3e-badge-color: var(--md-sys-color-on-primary);
}

.engine-badge[data-tone='tertiary'] {
  --m3e-badge-container-color: var(--md-sys-color-tertiary);
  --m3e-badge-color: var(--md-sys-color-on-tertiary);
}

.engine-badge[data-tone='error'] {
  --m3e-badge-container-color: var(--md-sys-color-error);
  --m3e-badge-color: var(--md-sys-color-on-error);
}
</style>
