<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { isAquesTalkVersion } from '../../backend/types';

const props = defineProps<{
  version: number;
}>();

const { t } = useI18n();

const safeVersion = computed(() => (isAquesTalkVersion(props.version) ? props.version : 1));

const badgeType = computed(() => {
  const map = {
    1: 'success',
    2: 'warning',
    10: 'danger',
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
  <var-badge
    v-bind="$attrs"
    class="engine-badge"
    :type="badgeType"
    :value="t(labelKey)"
  />
</template>

<style scoped>
.engine-badge {
  --badge-content-font-size: 0.75rem;
  --badge-content-padding: 2px 6px;
  box-shadow: 0 1px 3px hsla(var(--hsl-text), 0.2);
}

.engine-badge :deep(.var-badge__content),
.engine-badge :deep(.var-badge--right-top) {
  position: static;
  top: auto;
  right: auto;
  transform: none;
}
</style>
