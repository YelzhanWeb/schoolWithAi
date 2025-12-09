import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: {
      "/v1": {
        target: "http://localhost:8080",
        changeOrigin: true,
        secure: false,
      },
    },
    allowedHosts: [
      "school.furia.sbs",
      "furia.sbs",
      "www.furia.sbs",
      "api.furia.sbs",
      "ml.furia.sbs",
    ],
  },
});
