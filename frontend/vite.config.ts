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
			port: parseInt(env.VITE_PORT || '5173', 10),
			proxy: {
				'/api': {
					target: 'http://localhost:5181',
					changeOrigin: true
				}
			},
			allowedHosts: ["topological-gilberte-gynomonoecious.ngrok-free.dev"]
		}
	};
});
