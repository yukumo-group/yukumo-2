<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue';
import { X } from '@lucide/vue';
import { motion, useSpring } from 'motion-v';
import { useI18n } from 'vue-i18n';
import { resizeImageToDataUrl } from '../../../composables/resizeImage';

interface WizardUploadFile {
  file?: File;
  name?: string;
  url?: string;
  cover?: string;
  state?: 'loading' | 'success' | 'error';
}

const MAX_TILT_DEG = 16;

const props = defineProps<{
  active: boolean;
}>();

const imageDataUrl = defineModel<string | null>({ required: true });

const { t } = useI18n();

const uploadFiles = ref<WizardUploadFile[]>([]);
const rotateX = useSpring(0, { stiffness: 260, damping: 24 });
const rotateY = useSpring(0, { stiffness: 260, damping: 24 });

let pointerListening = false;

watch(
  imageDataUrl,
  (url) => {
    const current = uploadFiles.value[0]?.url;
    if (url === current) return;
    uploadFiles.value = url
      ? [{ url, cover: url, name: 'avatar.jpg', state: 'success' }]
      : [];
  },
  { immediate: true },
);

watch(
  () => props.active,
  (active) => {
    if (active) {
      startPointerTracking();
      return;
    }
    stopPointerTracking();
    rotateX.set(0);
    rotateY.set(0);
  },
  { immediate: true },
);

onBeforeUnmount(() => {
  stopPointerTracking();
});

async function onAfterRead(
  file: WizardUploadFile | WizardUploadFile[],
): Promise<void> {
  const item = Array.isArray(file) ? file[0] : file;
  if (!item?.file) return;
  try {
    item.state = 'loading';
    const dataUrl = await resizeImageToDataUrl(item.file);
    imageDataUrl.value = dataUrl;
    item.url = dataUrl;
    item.cover = dataUrl;
    item.state = 'success';
  } catch {
    item.state = 'error';
  }
}

function onRemoveFile(): void {
  imageDataUrl.value = null;
  uploadFiles.value = [];
}

function startPointerTracking(): void {
  if (pointerListening) return;
  pointerListening = true;
  window.addEventListener('pointermove', onDialogPointerMove);
}

function stopPointerTracking(): void {
  if (!pointerListening) return;
  pointerListening = false;
  window.removeEventListener('pointermove', onDialogPointerMove);
}

function onDialogPointerMove(event: PointerEvent): void {
  const dialog = document.querySelector('.character-wizard-dialog');
  if (!(dialog instanceof HTMLElement)) return;
  const rect = dialog.getBoundingClientRect();
  const inside =
    event.clientX >= rect.left &&
    event.clientX <= rect.right &&
    event.clientY >= rect.top &&
    event.clientY <= rect.bottom;
  if (!inside) {
    rotateX.set(0);
    rotateY.set(0);
    return;
  }
  const nx = (event.clientX - rect.left) / rect.width - 0.5;
  const ny = (event.clientY - rect.top) / rect.height - 0.5;
  rotateY.set(nx * 2 * MAX_TILT_DEG);
  rotateX.set(-ny * 2 * MAX_TILT_DEG);
}
</script>

<template>
  <div class="flex h-full flex-col items-center justify-center gap-3">
    <div class="uploader-stage">
      <motion.div
        class="uploader-tilt origin-center will-change-transform"
        :style="{ rotateX, rotateY }"
      >
        <var-uploader
          v-model="uploadFiles"
          class="wizard-uploader"
          accept="image/*"
          :maxlength="1"
          @after-read="onAfterRead"
          @remove="onRemoveFile"
        />
        <button
          v-if="imageDataUrl"
          type="button"
          class="uploader-remove"
          :aria-label="t('pages.profiles.wizard.removeImage')"
          v-ripple="{ color: 'rgba(255, 255, 255, 0.4)' }"
          @click.stop="onRemoveFile"
        >
          <X
            :size="16"
            aria-hidden="true"
          />
        </button>
      </motion.div>
    </div>
    <div class="flex flex-col items-center gap-1 text-center">
      <p class="text-sm text-on-surface-variant">
        {{ t('pages.profiles.wizard.imageHint') }}
      </p>
    </div>
  </div>
</template>

<style scoped>
.uploader-stage {
  perspective: 800px;
}

.uploader-tilt {
  position: relative;
}

.wizard-uploader {
  width: 11rem;
  --uploader-file-size: 11rem;
  --uploader-action-icon-size: 2.5rem;
  --uploader-file-margin: 0;
  --uploader-action-margin: 0;
  --uploader-file-border-radius: 1rem;
}

.wizard-uploader :deep(.var-uploader__file-list) {
  justify-content: center;
}

.wizard-uploader :deep(.var-uploader__file-close) {
  display: none;
}

.wizard-uploader :deep(.var-uploader__file-indicator) {
  display: none;
}

.wizard-uploader :deep(.var-form-details) {
  display: none;
}

.uploader-remove {
  position: absolute;
  top: 0.45rem;
  right: 0.45rem;
  z-index: 2;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.75rem;
  height: 1.75rem;
  padding: 0;
  border: none;
  border-radius: 9999px;
  background: rgba(0, 0, 0, 0.55);
  color: #fff;
  cursor: pointer;
  overflow: hidden;
  transition: background-color 160ms ease;
}

.uploader-remove:hover {
  background: rgba(0, 0, 0, 0.75);
}

.uploader-remove:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}
</style>
