# LOKARI Frontend (SvelteKit)

Frontend WebGIS LOKARI: SvelteKit 2 + Svelte 5 + Tailwind v4 + Leaflet —
6 halaman (`/`, `/danger-map`, `/safe-routes`, `/news`, `/search`, `/about`),
i18n ID/EN, mode terang/gelap, Web Push.

Panduan menjalankan, env, dan deploy (Vercel): [`README.md` di root](../README.md).
Detail arsitektur, pembagian `/api/*`, dan kontrak endpoint: [`docs/`](../docs/README.md).

```bash
npm install
npm run dev      # vite dev (+ proxy /api ke backend, kecuali /api/route milik SvelteKit)
npm run check    # svelte-check, wajib lolos sebelum merge
npm run build    # build produksi (adapter-vercel, runtime nodejs22.x)
```

Jangan tambahkan `/api/route` ke proxy Vite dan jangan teruskan ke backend
Go — endpoint itu ditangani SvelteKit sendiri (`src/routes/api/route/+server.ts`).
