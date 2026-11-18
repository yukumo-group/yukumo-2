import { createApp } from 'vue';
import { install as VueMonacoEditorPlugin } from '@guolao/vue-monaco-editor';

import App from './App.vue';
import router from './router';
import { i18n, initI18n } from './i18n';
import { initTheme } from './theme/theme';
import './m3e';
import './monaco/setup';
import './style.css';

initTheme();
initI18n();
const app = createApp(App);
app.config.compilerOptions.isCustomElement = (tag) => tag.startsWith('m3e-');
app.use(i18n).use(router).use(VueMonacoEditorPlugin).mount('#app');
