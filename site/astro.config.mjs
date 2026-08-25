import { defineConfig } from 'astro/config';

// Zerado landing page — fully static, zero client-side JavaScript.
// No integrations are added on purpose: the blueprint (design/blueprint.md §1.2)
// commits to zero JS, zero raster images and zero external requests.
export default defineConfig({
  site: 'https://zerado.app',
  output: 'static',
  compressHTML: true,
  build: {
    inlineStylesheets: 'auto'
  }
});
