import { ref } from 'vue';
import { getBackend } from '../backend';
import { normalizeVoice } from '../backend/params';
import type { Character } from '../backend/types';

const characters = ref<Character[]>([]);
const loading = ref(false);
const loaded = ref(false);
const error = ref('');
let pending: Promise<void> | null = null;

function normalizeCharacter(character: Character): Character {
  return {
    name: character.name.trim(),
    description: character.description.trim(),
    voices: (character.voices ?? []).map((voice) => normalizeVoice(voice)),
  };
}

async function refreshCharacters(): Promise<void> {
  loading.value = true;
  error.value = '';
  try {
    const list = await getBackend().listCharacters();
    characters.value = list.map((character) => ({
      name: character.name,
      description: character.description ?? '',
      voices: (character.voices ?? []).map((voice) => ({
        name: voice.name,
        params: voice.params ?? [],
      })),
    }));
    loaded.value = true;
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : String(cause);
    throw cause;
  } finally {
    loading.value = false;
  }
}

async function ensureLoaded(): Promise<void> {
  if (loaded.value) return;
  if (pending === null) {
    pending = refreshCharacters().finally(() => {
      pending = null;
    });
  }
  await pending;
}

async function addCharacter(character: Character): Promise<void> {
  await getBackend().addCharacter(normalizeCharacter(character));
  await refreshCharacters();
}

async function updateCharacter(
  oldName: string,
  character: Character,
): Promise<void> {
  await getBackend().updateCharacter(oldName, normalizeCharacter(character));
  await refreshCharacters();
}

async function deleteCharacter(name: string): Promise<void> {
  await getBackend().deleteCharacter(name);
  await refreshCharacters();
}

export function useCharacters() {
  return {
    characters,
    loading,
    loaded,
    error,
    refresh: refreshCharacters,
    ensureLoaded,
    addCharacter,
    updateCharacter,
    deleteCharacter,
  };
}
