import { defineConfig } from "vite";

export default defineConfig({
  build: {
    target: "es2022",
    outDir: "dist",
    emptyOutDir: true,

    lib: {
      entry: "src/content.ts",
      name: "NugsChat",
      formats: ["iife"],
      fileName: () => "content.js",
      cssFileName: 'styles'
    },

    // Useful while developing.
    sourcemap: true,
    minify: false,
  },
});
