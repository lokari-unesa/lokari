import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd(), '');
	return {
		plugins: [
			tailwindcss(),
			sveltekit()
		],
		server: {
			// Docker: VITE_PORT dari env container. Manual: dari file frontend/.env.
			port: parseInt(env.VITE_PORT || process.env.VITE_PORT || '5173', 10),

			// ============================================================
			// PROXY CONFIGURATION (Hanya aktif di mode DEV: `npm run dev`)
			// ============================================================
			// Mode Dev   : Proxy AKTIF — meneruskan request ke backend Go.
			// Mode Docker: Proxy TIDAK AKTIF — sudah dihandle oleh Docker network.
			//
			// PENTING: Endpoint /api/route TIDAK boleh di-proxy!
			//          Karena /api/route adalah SvelteKit server route
			//          (src/routes/api/route/+server.ts) untuk routing OSRM/ORS.
			//          Jika di-proxy, rute jalan akan menjadi garis lurus!
			//
			// Cara pakai:
			//   - Mode Dev (manual)  → BIARKAN proxy di bawah AKTIF (uncommented)
			//   - Mode Docker/Build  → Proxy otomatis TIDAK AKTIF, tidak perlu diubah
			// ============================================================
			proxy: {
				'/api/potensi': {
					target: process.env.BACKEND_URL || env.BACKEND_URL || 'http://localhost:5181',
					changeOrigin: true
				},
				'/api/news': {
					target: process.env.BACKEND_URL || env.BACKEND_URL || 'http://localhost:5181',
					changeOrigin: true
				},
				'/api/ai': {
					target: process.env.BACKEND_URL || env.BACKEND_URL || 'http://localhost:5181',
					changeOrigin: true
				},
				'/api/alert': {
					target: process.env.BACKEND_URL || env.BACKEND_URL || 'http://localhost:5181',
					changeOrigin: true
				},
				'/api/search': {
					target: process.env.BACKEND_URL || env.BACKEND_URL || 'http://localhost:5181',
					changeOrigin: true
				},
				'/api/health': {
					target: process.env.BACKEND_URL || env.BACKEND_URL || 'http://localhost:5181',
					changeOrigin: true
				}
			},

			allowedHosts: ["topological-gilberte-gynomonoecious.ngrok-free.dev"]
		}
	};
});
