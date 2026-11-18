import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import VueRouter from 'vue-router/vite';
import tailwindcss from '@tailwindcss/vite';
import wails from '@wailsio/runtime/plugins/vite';

const frontendRoot = path.dirname(fileURLToPath(import.meta.url));
const yukumoJsRoot = path.resolve(frontendRoot, '../yukumo-js');

const sharedArrayBufferHeaders = {
  'Cross-Origin-Opener-Policy': 'same-origin',
  'Cross-Origin-Embedder-Policy': 'require-corp',
};

export default defineConfig({
  server: {
    host: '127.0.0.1',
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
    headers: sharedArrayBufferHeaders,
    fs: {
      allow: [path.resolve(frontendRoot, '..')],
    },
    watch: {
      ignored: ['**/node_modules/.vite/**'],
    },
  },
  preview: {
    headers: sharedArrayBufferHeaders,
  },
  worker: {
    format: 'es',
  },
  assetsInclude: ['**/*.wasm'],
  optimizeDeps: {
    exclude: ['js7z-tools', 'yukumo.js', 'v86'],
  },
  resolve: {
    alias: [
      {
        find: 'yukumo.js/lang/zh',
        replacement: path.join(yukumoJsRoot, 'src/lang/zh/index.ts'),
      },
      {
        find: 'yukumo.js/lang/kanji2koe',
        replacement: path.join(yukumoJsRoot, 'src/lang/kanji2koe/index.ts'),
      },
      {
        find: 'yukumo.js/lang',
        replacement: path.join(yukumoJsRoot, 'src/lang/index.ts'),
      },
      {
        find: 'yukumo.js',
        replacement: path.join(yukumoJsRoot, 'src/index.ts'),
      },
      {
        find: 'v86',
        replacement: path.join(yukumoJsRoot, 'node_modules/v86'),
      },
    ],
  },
  plugins: [
    VueRouter({
      routesFolder: 'src/pages',
    }),
    vue(),
    tailwindcss(),
    wails('./bindings'),
  ],
});
