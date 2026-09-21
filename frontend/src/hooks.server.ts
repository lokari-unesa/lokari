import type { Handle } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

// Endpoint /api yang benar-benar di-handle SvelteKit sendiri (bukan backend Go)
const KIT_API_PREFIXES = ['/api/route'];

// BACKEND_URL diisi lewat compose (base) untuk mode production;
// di mode dev, vite proxy yang menangani /api (hooks ini nyaris tidak terpanggil).
const BACKEND_URL = env.BACKEND_URL ?? 'http://localhost:5181';

export const handle: Handle = async ({ event, resolve }) => {
	const { url, request } = event;
	const path = url.pathname;

	// Mode build/prod tidak punya vite proxy, jadi /api/* diteruskan ke backend Go
	if (path.startsWith('/api/') && !KIT_API_PREFIXES.some((p) => path.startsWith(p))) {
		try {
			const upstream = await fetch(new URL(path + url.search, BACKEND_URL), {
				method: request.method,
				headers: {
					'content-type': request.headers.get('content-type') ?? '',
					accept: request.headers.get('accept') ?? 'application/json'
				},
				body: ['GET', 'HEAD'].includes(request.method) ? undefined : await request.text()
			});
			return new Response(await upstream.text(), {
				status: upstream.status,
				headers: { 'content-type': upstream.headers.get('content-type') ?? 'application/json' }
			});
		} catch {
			return new Response(JSON.stringify({ error: 'backend unreachable' }), {
				status: 502,
				headers: { 'content-type': 'application/json' }
			});
		}
	}

	return resolve(event);
};