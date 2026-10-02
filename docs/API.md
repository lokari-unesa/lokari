# Kontrak REST API — LOKARI

> Sumber kebenaran: `backend/internal/api/routes.go` (+ handler di
> `backend/internal/api/handlers/`) dan
> `frontend/src/routes/api/route/+server.ts`.
>
> Semua path di bawah ber-prefix `/api` (bukan `/api/v1` seperti proposal lama).
> Terakhir diselaraskan dengan kode: Oktober 2026.

## 1. Aturan routing `/api/*` (wajib paham sebelum ubah apa pun)

| Lingkungan | Mekanisme | Keterangan |
|---|---|---|
| Dev manual (`npm run dev`) | Proxy Vite per-path eksplisit (`frontend/vite.config.ts`) | Hanya path yang terdaftar yang diteruskan ke `BACKEND_URL` (default `http://localhost:5181`) |
| Docker / build / Vercel | `frontend/src/hooks.server.ts` meneruskan ke `BACKEND_URL` server-side | Tanpa masalah CORS |

**Pengecualian: `GET /api/route` ditangani SvelteKit sendiri, bukan backend
Go** (`KIT_API_PREFIXES = ['/api/route']`). Jangan masukkan `/api/route` ke
proxy Vite dan jangan teruskan ke backend — jika diteruskan, respons rute
jalan jatuh ke garis lurus (straight line) karena backend tidak mengenal
endpoint itu.

## 2. Endpoint backend Go (Fiber, port 5181)

| Method & Path | Handler | Rate limit | Respons |
|---|---|---|---|
| `GET /` | — | — | String sapaan (bukan 404) |
| `GET /api/health` | — | — | `{status, message}` — dipakai healthcheck Docker |
| `GET /api/potensi` | `Handler.GetPotensiBencana` | — | `{status: "success", data: [...]}` — `geometri` berupa string GeoJSON; baris tanpa geometri dilewati + dilog |
| `POST /api/search` | `AIHandler.SearchSemantic` | 30 req/menit/IP | `{message, data: [...15 hasil]}` — body `{query, lat?, lng?}`; query maks 500 karakter; `embedding` via Cohere 1024-d, `ORDER BY embedding <=> $1`; `deskripsi`/`kapasitas` NULL dikirim `null` |
| `GET /api/alert` | `AIHandler.GetAlert` | 30 req/menit/IP | ⚠️ **FITUR MATI** — endpoint hidup, tetapi tidak pernah dipanggil frontend mana pun |
| `GET /api/news` | `NewsHandler.GetNews` | — | `{status: "success", data: [...maks 20, terbaru dulu]}` — field `created_at` RFC3339 zona Asia/Jakarta |
| `GET /api/kelud/status` | `StatusHandler.GetKeludStatus` | — | `{status: "success", data: {status: "Level II (Waspada)", updated_at}}` atau `data: null` bila belum ada state — satu-satunya sumber kebenaran kartu status di beranda |
| `POST /api/subscribe` | `PushHandler.Subscribe` | — | Mendaftarkan Web Push subscription browser ke `push_subscriptions` |

Batas global: `BodyLimit` 64 KB (hanya `POST /api/search` dan
`POST /api/subscribe` yang memakai body). Error internal tidak pernah
dibocorkan ke client — handler mengembalikan pesan generik
(`serverError`) dan menulis detail ke log server. CORS whitelist default
`http://localhost:5180,http://localhost:5173`, override via env
`CORS_ORIGINS`.

## 3. Endpoint SvelteKit (`GET /api/route`)

Query: `?start=lng,lat&end=lng,lat` — 4 angka wajib finite dalam rentang
geografis, kalau tidak balas 400.

1. **Strategi 1 — OpenRouteService**: `POST
   https://api.openrouteservice.org/v2/directions/driving-car/geojson`
   dengan `avoid_polygons` (poligon zona bahaya hardcoded di file yang sama).
   Butuh `ORS_API_KEY` (server-only). Hasil bertanda
   `{isSafe: true, status: "safe"}`.
2. **Strategi 2 — fallback OSRM**: `router.project-osrm.org`, tanpa
   penghindaran bahaya. Hasil bertanda `{isSafe: false, status: "fallback"}`.
3. Keduanya gagal → 502. Frontend menampilkan 3 state (aman / fallback /
   bahaya) dari `isSafe` + `status`; tanpa ETA dari routing engine, ETA
   tampil "—".

## 4. Env yang dibaca endpoint

| Variabel | Dibaca di | Untuk |
|---|---|---|
| `ORS_API_KEY` | `api/route/+server.ts` | OpenRouteService (wajib untuk rute aman) |
| `BACKEND_URL` | `hooks.server.ts`, proxy Vite | Target penerusan `/api/*` |
| `AI_API_KEY` / `AI_BASE_URL` / `AI_MODEL` | `internal/ai/nlp.go` | Provider OpenAI-compatible (ringkasan berita) |
| `COHERE_API_KEY` | `internal/ai/nlp.go` | Embedding `embed-multilingual-v3.0` 1024-d |
| `VAPID_PUBLIC_KEY` / `VAPID_PRIVATE_KEY` | worker `broadcastPush` | Web Push (urgency high, TTL 1 jam) |

Skema tabel di balik endpoint: `docs/Schema.md`. Alur worker yang mengisi
`kabar_kelud`/`monitor_state`: `docs/Architecture.md`.
