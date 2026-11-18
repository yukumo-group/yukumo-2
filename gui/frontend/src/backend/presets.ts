import { ref } from 'vue';
import { getBackend } from './index';

const namesByVersion = ref<Record<number, string[]>>({});
const pending = new Map<number, Promise<string[]>>();

async function loadPresetNames(version: number): Promise<string[]> {
  const cached = namesByVersion.value[version];
  if (cached) return cached;
  const existing = pending.get(version);
  if (existing) return existing;

  const request = getBackend()
    .presetNames(version)
    .then((names) => {
      namesByVersion.value = { ...namesByVersion.value, [version]: names };
      return names;
    })
    .finally(() => {
      pending.delete(version);
    });
  pending.set(version, request);
  return request;
}

export function usePresetCatalog() {
  async function ensure(version: number): Promise<string[]> {
    return loadPresetNames(version);
  }

  function presetNames(version: number): string[] {
    return namesByVersion.value[version] ?? [];
  }

  function presetLabel(version: number, preset: number): string {
    return presetNames(version)[preset] ?? `#${preset}`;
  }

  return { ensure, presetNames, presetLabel };
}
