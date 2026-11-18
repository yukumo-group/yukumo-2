<template>
  <m3e-card class="settings-card overflow-hidden rounded-2xl" variant="elevated">
    <slot />
  </m3e-card>
</template>

<style scoped>
:deep(.settings-item + .settings-item) {
  position: relative;
}

:deep(.settings-item + .settings-item)::before {
  content: '';
  position: absolute;
  top: 0;
  right: 1rem;
  left: 1rem;
  border-top: 1px solid var(--md-sys-color-on-surface-variant);
  opacity: 0.2;
  pointer-events: none;
}
</style>
