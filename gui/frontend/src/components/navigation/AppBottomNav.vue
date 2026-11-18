<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';
import { bottomDestinations } from '../../navigation/destinations';
import NavIcon from './NavIcon.vue';
import NavLabel from './NavLabel.vue';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();

const active = computed(() =>
  route.path === '/help' ? '/settings' : route.path,
);

function onSelect(path: string): void {
  if (path !== route.path) {
    void router.push(path);
  }
}
</script>

<template>
  <nav
    class="relative z-20 shrink-0 sm:hidden"
    :aria-label="t('nav.mainAriaLabel')"
  >
    <m3e-nav-bar>
      <m3e-nav-item
        v-for="destination in bottomDestinations"
        :key="destination.path"
        :href="`#${destination.path}`"
        .selected="active === destination.path"
        @click.prevent="onSelect(destination.path)"
      >
        <NavIcon
          slot="icon"
          :icon="destination.icon"
          :active="active === destination.path"
          :variant="destination.path === '/settings' ? 'settings' : 'default'"
        />
        <NavLabel
          :active="active === destination.path"
          :text="t(destination.labelKey)"
          :gap="6"
        />
      </m3e-nav-item>
    </m3e-nav-bar>
  </nav>
</template>
