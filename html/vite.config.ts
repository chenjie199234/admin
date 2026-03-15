import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@api': path.resolve(__dirname, '../api') // 根据实际 api 目录调整
    }
  },
  server:{
	host:'0.0.0.0',
  }
})
