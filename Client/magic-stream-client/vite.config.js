import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    watch: {
      usePolling: process.env.DOCKER_WATCH_POLLING === 'true',
      interval: 500,
    },
  },
})
