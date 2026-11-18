<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';
import {
  helpDestination,
  primaryDestinations,
  settingsDestination,
} from '../../navigation/destinations';
import Hint from '../ui/Hint.vue';
import NavIcon from './NavIcon.vue';
import NavLabel from './NavLabel.vue';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();

const mainDestinations = [...primaryDestinations, helpDestination];

const active = computed(() => route.path);

function hrefFor(path: string): string {
  return `#${path}`;
}

function onSelect(path: string): void {
  if (path !== route.path) {
    void router.push(path);
  }
}
</script>

<template>
  <aside
    class="rail-chrome relative z-20 hidden h-full shrink-0 sm:flex flex-col"
    :aria-label="t('nav.mainAriaLabel')"
  >
    <m3e-nav-rail class="h-full">
      <Hint
        v-for="destination in mainDestinations"
        :key="destination.path"
        :text="t(destination.labelKey)"
        position="after"
        :disabled="active === destination.path"
      >
        <template #default="{ id }">
          <m3e-nav-item
            :id="id"
            :href="hrefFor(destination.path)"
            .selected="active === destination.path"
            @click.prevent="onSelect(destination.path)"
          >
            <NavIcon
              slot="icon"
              :icon="destination.icon"
              :active="active === destination.path"
            />
            <NavLabel
              :active="active === destination.path"
              :text="t(destination.labelKey)"
            />
          </m3e-nav-item>
        </template>
      </Hint>
      <div class="rail-spacer" aria-hidden="true" />
      <m3e-nav-item
        :href="hrefFor(settingsDestination.path)"
        .selected="active === settingsDestination.path"
        @click.prevent="onSelect(settingsDestination.path)"
      >
        <NavIcon
          slot="icon"
          :icon="settingsDestination.icon"
          :active="active === settingsDestination.path"
          variant="settings"
        />
        <NavLabel
          :active="active === settingsDestination.path"
          :text="t(settingsDestination.labelKey)"
        />
      </m3e-nav-item>
    </m3e-nav-rail>
  </aside>
</template>

<style scoped>
.rail-chrome {
  background-color: var(--md-sys-color-surface-container-low);
}

.rail-chrome m3e-nav-rail {
  /* Tailwind's margin reset cancels the rail's negative item margin, so the
     default inline padding shifts icons left and clips them on the right. */
  --m3e-nav-rail-inline-padding: 0px;
}

.rail-spacer {
  flex: 1 1 auto;
  align-self: stretch;
  min-height: 0.75rem;
}
</style>
