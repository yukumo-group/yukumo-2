<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Plus } from '@lucide/vue';
import { Snackbar } from '@varlet/ui';
import { useI18n } from 'vue-i18n';
import CharacterCard from '../components/profiles/CharacterCard.vue';
import CharacterWizardDialog from '../components/profiles/CharacterWizardDialog.vue';
import type { Character } from '../backend/types';
import { useCharacters } from '../store/useCharacters';

const { t } = useI18n();
const { characters, loading, error, ensureLoaded } = useCharacters();

const wizardOpen = ref(false);
const editing = ref<Character | null>(null);

onMounted(() => {
  void ensureLoaded().catch(reportError);
});

function reportError(cause: unknown): void {
  Snackbar({
    content: cause instanceof Error ? cause.message : String(cause),
    type: 'error',
  });
}

function openCreate(): void {
  editing.value = null;
  wizardOpen.value = true;
}

function openEdit(character: Character): void {
  editing.value = character;
  wizardOpen.value = true;
}
</script>

<template>
  <section
    class="relative flex w-full max-w-3xl flex-col gap-6 px-6 pt-8 pb-8 text-left sm:px-6"
  >
    <header class="flex items-center justify-between gap-4">
      <h1 class="min-w-0 font-bold tracking-wide text-text text-3xl">
        {{ t('pages.profiles.title') }}
      </h1>
      <var-button
        class="shrink-0"
        type="primary"
        :aria-label="t('pages.profiles.createAriaLabel')"
        @click="openCreate"
      >
        <Plus :size="18" aria-hidden="true" class="mr-2" />
        {{ t('pages.profiles.create') }}
      </var-button>
    </header>

    <p
      v-if="error"
      class="rounded-lg bg-danger/15 px-4 py-3 text-sm text-danger"
      role="alert"
    >
      {{ error }}
    </p>

    <div
      v-if="!loading && characters.length === 0"
      class="flex flex-col items-start gap-3 rounded-2xl bg-surface-container px-4 py-6"
    >
      <p class="text-text">{{ t('pages.profiles.empty') }}</p>
      <p class="text-sm text-on-surface-variant">
        {{ t('pages.profiles.emptyHint') }}
      </p>
    </div>

    <div
      v-else
      class="character-grid"
    >
      <CharacterCard
        v-for="character in characters"
        :key="character.name"
        :character="character"
        @select="openEdit(character)"
      />
    </div>

    <CharacterWizardDialog
      v-model:show="wizardOpen"
      :character="editing"
    />
  </section>
</template>

<style scoped>
.character-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(10rem, 1fr));
  gap: 1rem;
}
</style>
