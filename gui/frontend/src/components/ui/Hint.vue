<script setup lang="ts">
import { ref, watch } from 'vue';
import { useId } from 'vue';
import type { M3eTooltipElement, TooltipPosition } from '@m3e/web/tooltip';

const props = defineProps<{
  text: string;
  position?: TooltipPosition;
  disabled?: boolean;
  forceOpen?: boolean;
}>();

const id = useId();
const tooltipRef = ref<M3eTooltipElement | null>(null);

watch(
  () => props.forceOpen,
  (open) => {
    if (open === undefined) return;
    const tooltip = tooltipRef.value;
    if (!tooltip) return;
    if (open) void tooltip.show();
    else tooltip.hide();
  },
);
</script>

<template>
  <span class="contents">
    <slot :id="id" />
    <m3e-tooltip
      ref="tooltipRef"
      :for="id"
      :position="position ?? 'below'"
      .disabled="disabled ?? false"
    >
      {{ text }}
    </m3e-tooltip>
  </span>
</template>
