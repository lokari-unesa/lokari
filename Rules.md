# Project Rules & Guidelines - LOKARI

> Aturan di bawah mengikuti **stack yang benar-benar dipakai**. Batasan yang
> terbukti tidak diimplementasikan (pgRouting, port 6000/6001, Nginx/PM2)
> sudah dihapus — alasannya dicatat di `General.md` §3 agar reviewer
> memahami konteksnya.

## 1. Batasan Stack

- **[DILARANG]** React, Next.js, Vue, atau Base44 SDK (`@base44/sdk` wajib
  dibuang bila memigrasi prototipe lama).
- **[DILARANG]** MongoDB atau NoSQL lainnya.
- **[WAJIB]** SvelteKit 2 + Svelte 5 (frontend), Golang + Fiber (backend),
  PostgreSQL 16 + PostGIS + pgvector (database).
- Routing evakuasi memakai OpenRouteService + fallback OSRM di endpoint
  SvelteKit (`src/routes/api/route/+server.ts`) — jangan membangun ulang
  dengan pgRouting tanpa digitasi topologi jalan terlebih dahulu.

## 2. Aturan Backend Golang (Zero-Admin & AI)

1. **Background Workers:** fetch BMKG/MAGMA/EONET berjalan di cron dual-loop
   (hot 30 dtk + cold 6 jam, mutex anti-tumpang-tindih). Jangan menunda
   respons API ke client gara-gara server sedang menarik data eksternal.
   Trigger manual: `go run scripts/ops.go hot|cold|all`.
2. **Koneksi Database:** driver `jackc/pgx` (pool). Semua handler harus guard
   `h.DB == nil` — server wajib tetap hidup tanpa DB.
3. **HTTP keluar:** selalu pakai `http.Client{Timeout: ...}` dan cek
   `resp.StatusCode`. AI wajib punya fallback deterministik + guard
   `len(resp.Choices) == 0`.
4. **Berita & notifikasi:** insert selalu `ON CONFLICT (sumber, judul) DO
   NOTHING`; broadcast push hanya untuk kategori penting; subscription mati
   (410/404) dihapus.
5. **Security & Configuration:** port aktual 5181 (backend), 5180/5173
   (frontend); DSN, `AI_*`, `COHERE_API_KEY`, `VAPID_*`, `ORS_API_KEY`
   **WAJIB** via `.env` — tidak boleh hardcode. Error internal dilog di
   server, client menerima pesan generik (`serverError`). Endpoint AI
   ber-rate-limit (30 req/menit/IP) + `BodyLimit` 64 KB.

## 3. Aturan Frontend SvelteKit

1. **Pemetaan:** Leaflet.js (dynamic import, aman SSR). Konversi GeoJSON
   dari Golang via helper toleran string/object (`$lib/geo.ts`) — jangan
   `JSON.parse` mentah.
2. **Routing `/api/*`:** `/api/route` milik SvelteKit dan **dilarang**
   masuk proxy Vite / diteruskan ke backend. Endpoint lain diteruskan via
   proxy Vite (dev) atau `hooks.server.ts` (prod). Lihat `docs/API.md` §1.
3. **State server:** kartu status Kelud hanya dari `GET /api/kelud/status`;
   tampilan rute wajib menghormati flag `isSafe`/`status` (3 state:
   aman/fallback/bahaya); tanpa ETA dari engine, tampilkan "—".
4. **Ringan & i18n:** aset WebP, string UI via store i18n (ID/EN) — jangan
   hardcode teks baru di komponen. `@inlang/paraglide-js` terpasang tapi
   tidak dipakai; jangan impor dari sana.

## 4. Infrastruktur & Testing

- **Deployment:** frontend via Vercel (`frontend/` sebagai Root Directory,
  runtime `nodejs22.x`, env `BACKEND_URL` + `ORS_API_KEY` via dashboard);
  backend + DB via Docker (image prod non-root + healthcheck
  `/api/health`). `frontend/Dockerfile` hanya untuk dev.
- **Database:** skrip idempoten (`migrate.go` → `seed.go`); `migrate.go
  reset` destruktif — dilarang di DB berisi data. Base image DB wajib
  `pgvector/pgvector:pg16` + PostGIS.
- **Verifikasi:** `go vet ./...` + `npm run check` + `npm run build`
  sebelum merge (CI hanya jalan di `master` repo `lokari-unesa/lokari`).
  Belum ada unit test (`TODO.md` 6.1) — pengecualian sementara, bukan standar.

---

*Status: SELARAS DENGAN IMPLEMENTASI (Okt 2026).*
