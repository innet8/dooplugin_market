import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

export default defineConfig({
  plugins: [react()],
  base: "/store",
  server: {
    host: '0.0.0.0',
    hmr:true,
    proxy: {
      '/api': {
        target: 'http://192.168.31.214:8080',
        changeOrigin: true,
      }
    },

    headers: {
      'Access-Control-Allow-Origin': '*',  // 允许跨域连接
    },
    
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@locales': '/web/public/locales'
    },
  },
});
