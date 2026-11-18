<script setup lang="ts">
import { useId } from 'vue';
import { useI18n } from 'vue-i18n';
import type { AquesTalkVersion } from '../../backend/types';
import EngineBadge from '../profiles/EngineBadge.vue';
import VersionSelectLabel from '../profiles/VersionSelectLabel.vue';

const ENGINE_VERSIONS: AquesTalkVersion[] = [1, 2, 10];

const versionLabelKey: Record<AquesTalkVersion, string> = {
  1: 'pages.profiles.engines.at1',
  2: 'pages.profiles.engines.at2',
  10: 'pages.profiles.engines.at10',
};

const props = defineProps<{
  modelValue: AquesTalkVersion;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: AquesTalkVersion];
}>();

const { t } = useI18n();
const triggerId = useId();
const menuId = useId();

function selectVersion(version: AquesTalkVersion): void {
  emit('update:modelValue', version);
}
</script>

<template>
  <span class="inline-flex shrink-0" @pointerdown.stop>
    <button
      :id="triggerId"
      type="button"
      class="inline-flex shrink-0 items-center rounded-sm hover:opacity-80"
      :aria-label="
        t('pages.generate.timeline.engineAriaLabel', {
          engine: t(versionLabelKey[props.modelValue]),
        })
      "
      aria-haspopup="menu"
    >
      <EngineBadge :version="modelValue" />
      <m3e-menu-trigger :for="menuId" />
    </button>
    <m3e-menu
      :id="menuId"
      class="timeline-engine-menu"
      position-x="before"
      position-y="below"
    >
      <m3e-menu-item
        v-for="version in ENGINE_VERSIONS"
        :key="version"
        @click="selectVersion(version)"
      >
        <VersionSelectLabel
          :version="version"
          :label="t(versionLabelKey[version])"
        />
      </m3e-menu-item>
    </m3e-menu>
  </span>
</template>

<style>
.timeline-engine-menu {
  overflow: hidden;
}
</style>
