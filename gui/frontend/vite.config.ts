import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import VueRouter from 'vue-router/vite';
import tailwindcss from '@tailwindcss/vite';
import wails from '@wailsio/runtime/plugins/vite';

export default defineConfig({
  server: {
    host: '127.0.0.1',
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
  assetsInclude: ['**/*.wasm'],
  plugins: [
    VueRouter({
      routesFolder: 'src/pages',
    }),
    vue({
      template: {
        compilerOptions: {
          isCustomElement: (tag) => tag.startsWith('m3e-'),
        },
      },
    }),
    tailwindcss(),
    wails('./bindings'),
  ],
});
