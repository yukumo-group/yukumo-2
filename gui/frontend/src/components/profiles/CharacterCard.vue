<script setup lang="ts">
import { computed } from 'vue';
import { User } from '@lucide/vue';
import { useI18n } from 'vue-i18n';
import type { Character } from '../../backend/types';
import { voiceVersions } from '../../backend/params';
import EngineBadge from './EngineBadge.vue';

const props = withDefaults(
  defineProps<{
    character: Pick<Character, 'name' | 'description' | 'voices'>;
    interactive?: boolean;
  }>(),
  {
    interactive: true,
  },
);

const emit = defineEmits<{
  select: [];
}>();

const { t } = useI18n();

const ariaLabel = computed(() =>
  t('pages.profiles.cardAriaLabel', { name: props.character.name }),
);

const engines = computed(() => {
  const seen = new Set<number>();
  const versions: number[] = [];
  for (const voice of props.character.voices) {
    for (const version of voiceVersions(voice)) {
      if (seen.has(version)) continue;
      seen.add(version);
      versions.push(version);
    }
  }
  return versions;
});
</script>

<template>
  <component
    :is="interactive ? 'button' : 'article'"
    class="character-card relative flex flex-col overflow-hidden rounded-2xl bg-surface-container text-left text-on-surface"
    :type="interactive ? 'button' : undefined"
    :aria-label="interactive ? ariaLabel : undefined"
    @click="interactive ? emit('select') : undefined"
  >
    <m3e-card
      class="flex h-full flex-col overflow-hidden rounded-2xl"
      variant="elevated"
    >
      <div
        class="character-thumb relative flex shrink-0 items-center justify-center bg-surface-container-high text-on-surface-variant"
      >
        <User :size="40" aria-hidden="true" />
        <div
          v-if="engines.length > 0"
          class="pointer-events-none absolute top-2 left-2 flex flex-wrap gap-1"
        >
          <EngineBadge
            v-for="version in engines"
            :key="version"
            :version="version"
          />
        </div>
      </div>

      <div class="flex flex-1 flex-col gap-1 p-3 pt-2">
        <h2 class="truncate font-bold text-on-surface text-base">
          {{ character.name }}
        </h2>
        <p
          v-if="character.description"
          class="line-clamp-2 text-sm leading-snug text-on-surface-variant"
        >
          {{ character.description }}
        </p>
        <p class="truncate text-xs text-on-surface-variant">
          {{ t('pages.profiles.voiceCount', { n: character.voices.length }) }}
        </p>
      </div>
    </m3e-card>
  </component>
</template>

<style scoped>
.character-card {
  cursor: default;
  width: 100%;
  min-width: 0;
}

.character-thumb {
  width: 100%;
  aspect-ratio: 1;
}

button.character-card {
  cursor: pointer;
}

button.character-card :deep(m3e-card)::after {
  content: '';
  position: absolute;
  inset: 0;
  z-index: 1;
  border-radius: inherit;
  pointer-events: none;
  background: color-mix(in srgb, var(--md-sys-color-on-surface) 8%, transparent);
  opacity: 0;
  transition: opacity 120ms ease;
}

button.character-card:hover :deep(m3e-card)::after {
  opacity: 1;
}

button.character-card:focus-visible :deep(m3e-card)::after,
button.character-card:active :deep(m3e-card)::after {
  background: color-mix(in srgb, var(--md-sys-color-on-surface) 12%, transparent);
  opacity: 1;
}

button.character-card:focus-visible {
  outline: 2px solid var(--md-sys-color-primary);
  outline-offset: 2px;
}
</style>
