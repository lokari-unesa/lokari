## Legenda Prioritas

| Label | Arti |
|-------|------|
| 🔴 **KRITIS** | Bug fungsional yang merusak fitur / masalah keamanan. Perbaiki duluan. |
| 🟠 **SEDANG** | Bug robusteness / risiko operasional yang perlu ditangani. |
| 🟡 **RENDAH** | Pelanggaran best practice / kebersihan kode. |
| ✅ | Sudah selesai dikerjakan |

---

## 1. 🔴 KRITIS — Bug Fungsional

- [ ] **1.1 Semua berita tampil sebagai kategori "volcano"**
  - **File:** `frontend/src/routes/news/+page.svelte:11-17, 37` dan `frontend/src/routes/+page.svelte:46-52, 77, 83`
  - **Masalah:** `categoryToSlug` hanya memetakan label Indonesia (`"Gunung Api"`, `"Lahar"`, dst.), sementara DB `kabar_kelud.kategori` berisi slug *lowercase English* hasil AI (`volcano`, `lahar`, `evac`, `weather`, `warning`). Karena `categoryToSlug[n.category]` selalu `undefined`, fallback `"volcano"` dipakai untuk SEMUA berita → ikon & filter di halaman News dan beranda salah total; filter "Lahar/Evakuasi/Cuaca/Peringatan" selalu kosong.
  - **Perbaikan:** `categoryToSlug` harus memetakan kedua bentuk (slug EN + label ID), atau backend mengembalikan slug EN langsung.

- [ ] **1.2 Crash `JSON.parse` pada marker di halaman `/safe-routes`**
  - **File:** `frontend/src/routes/safe-routes/+page.svelte:50-71` (crash di baris 58 & 63); data di-parse di `:123-129`
  - **Masalah:** `fetchShelters()` sudah mengubah `geometri` menjadi *object* (via `JSON.parse`), lalu `onMarkerClick` memanggil `JSON.parse(potensi.geometri)` lagi — `JSON.parse(object)` melempar `SyntaxError`. Popup tetap terbuka tapi rekomendasi faskes tidak pernah muncul + error di console. (`/danger-map` aman karena geometrinya masih string.)
  - **Perbaikan:** Buat satu helper `getGeoCoordinates(obj)` yang toleran terhadap string/object, pakai di semua halaman.

- [ ] **1.3 Rute `/api/route` mati di dev; seluruh `/api/*` mati di produksi**
  - **File:** `frontend/vite.config.ts:14-19`, `frontend/src/routes/api/route/+server.ts`, `frontend/svelte.config.js:1,9`
  - **Masalah:** Di dev, proxy Vite menangkap SEMUA `/api` → backend Fiber, termasuk `/api/route` yang harus di-handle endpoint SvelteKit → 404 → rute jatuh diam-diam ke garis lurus. Di produksi (`adapter-auto`), tidak ada proxy dan `/api/potensi`, `/api/news`, `/api/search` tidak punya endpoint SvelteKit → semua 404. Aplikasi hanya berfungsi di setup Docker dev-mode.
  - **Perbaikan:** Pilih satu arsitektur — (a) semua `/api` dari backend + hilangkan endpoint SvelteKit `/api/route` (taruh routing di backend), atau (b) pasang `adapter-node` + reverse proxy Nginx `/api` → backend untuk produksi.

- [ ] **1.4 Klaim "Aman digunakan" selalu hijau walau rute tidak aman**
  - **File:** `frontend/src/routes/safe-routes/+page.svelte:545-550` (UI), `frontend/src/routes/api/route/+server.ts:63,82` (flag `isSafe`), `frontend/src/lib/components/KeludMapView.svelte:100-118` (`fetchRealRoute` mengabaikan `isSafe`)
  - **Masalah:** UI selalu menampilkan `CheckCircle2` + "Aman digunakan" meskipun server mengembalikan rute fallback OSRM dengan `isSafe: false`. ETA (`dist/15 km/h`) dan "Rute Alternatif" (`distanceKm * 1.5`, baris 591/595) adalah angka karangan, bukan dari routing engine. Untuk aplikasi penyelamatan, ini menyesatkan dan berbahaya.
  - **Perbaikan:** Teruskan flag `isSafe` ke UI, tampilkan peringatan bila `false`/fallback, hapus metrik buatan. Validasi numerik input rute (`isNaN`/`isFinite`) di `api/route/+server.ts:11-14`.

- [ ] **1.5 Hasil semantic search diam-diam berkurang (NULL di-skip)**
  - **File:** `backend/internal/api/handlers/ai_handler.go:110-114`
  - **Masalah:** Kolom `deskripsi` & `kapasitas` di-scan ke `string`/`int` (bukan pointer). Jika NULL di DB, `rows.Scan` error → `continue` → baris di-skip tanpa log, hasil < 15 item secara misterius.
  - **Perbaikan:** Scan ke `*string`/`*int` atau `pgtype`, lalu log count baris yang di-skip; tambahkan `rows.Err()` setelah loop.

- [ ] **1.6 Timestamp berita dibuang — semua tampil "Baru Saja"**
  - **File:** `backend/internal/api/handlers/news_handler.go:46-56`
  - **Masalah:** `created_at` di-scan ke `interface{}` lalu diabaikan dan diganti string hardcoded `"Baru Saja"`.
  - **Perbaikan:** Format `CreatedAt` dengan layout `time.RFC3339`/lokal Indonesia dan kirim ke frontend.

## 2. 🔴 KRITIS — Keamanan

- [ ] **2.1 API key OpenRouteService di-hardcode di source code**
  - **File:** `frontend/src/routes/api/route/+server.ts:46`
  - **Masalah:** Key `5b3ce359[...]` tertanam langsung di kode. Melanggar `Rules.md §2.4` ("kunci API eksternal WAJIB disimpan di .env, tidak boleh hardcode"). Key ini juga ter-commit ke git history.
  - **Perbaikan:** Pindah ke `env` (server-side), tambahkan `.env.example`, dan **rotate/revoke** key yang sudah terlanjur bocor di repo.

- [ ] **2.2 Endpoint AI tanpa otentikasi & rate limit**
  - **File:** `backend/internal/api/handlers/ai_handler.go` (`/api/search` & `/api/alert`), `backend/internal/api/routes.go:32-33`
  - **Masalah:** `/api/search` memakai kuota Cohere berbayar per request. Siapa pun bisa memanggil terus-menerus dan menguras kuota (DoS terhadap anggaran). `POST /api/search` juga tidak membatasi ukuran body.
  - **Perbaikan:** Tambah rate limiting (Fiber limiter), cap ukuran body, dan validasi server-side isi query (jangan andalkan blocklist client).

- [ ] **2.3 Error internal/basis data bocor ke client**
  - **File:** `backend/internal/api/handlers/handlers.go:33,46`, `backend/internal/api/handlers/ai_handler.go:104`
  - **Masalah:** `err.Error()` dari DB/SQL dikembalikan mentah sebagai response JSON → membocorkan detail skema/internal.
  - **Perbaikan:** Log error di server, kembalikan pesan generik ke client.

- [ ] **2.4 Kredensial & konfigurasi DB di-hardcode di compose**
  - **File:** `docker-compose.yaml:35,46` (`postgres:root`)
  - **Masalah:** Password `root` tertulis jelas; port DB (5433) dibuka ke host.
  - **Perbaikan:** Pakai env dengan nilai default yang kuat/random, batasi expose port DB hanya untuk dev lokal.

- [ ] **2.5 CORS terbuka penuh (`*`)**
  - **File:** `backend/cmd/server/main.go:45-48`
  - **Masalah:** `AllowOrigins: "*"` + tidak ada channel otentikasi → endpoint apa pun bisa dipanggil dari situs lain.
  - **Perbaikan:** Whitelist origin frontend yang dikenal; buat daftar `AllowMethods/AllowHeaders` eksplisit.

## 3. 🟠 SEDANG — Robustness Backend

- [ ] **3.1 Request HTTP keluar tanpa timeout**
  - **File:** `backend/internal/api/handlers/ai_handler.go:29` (`http.Get` NASA), `backend/internal/service/fetcher.go:84` (`http.Get` NASA), `backend/internal/ai/nlp.go:150` (`http.Client{}` kosong untuk Cohere)
  - **Masalah:** Tidak ada timeout → handler/worker bisa menggantung. Di `GetAlert`, status code respons NASA juga tidak dicek.
  - **Perbaikan:** Pakai `http.Client{Timeout: ...}` dan cek `resp.StatusCode`; beri fallback yang jelas.

- [ ] **3.2 `h.DB == nil` → panic di handler AI & News**
  - **File:** `backend/internal/api/handlers/ai_handler.go`, `news_handler.go`
  - **Masalah:** Hanya `GetPotensiBencana` yang guard `h.DB == nil`. Jika `DATABASE_URL` kosong di `.env`, handler AI/News memanggil method pada pool nil → panic (ditangkap recover) → 500 kosong.
  - **Perbaikan:** Guard konsisten di semua handler (atau reject startup bila DB wajib).

- [ ] **3.3 `rows.Err()` tidak pernah dicek**
  - **File:** `backend/internal/api/handlers/handlers.go`, `ai_handler.go`, `news_handler.go`, `backend/internal/service/fetcher.go`
  - **Masalah:** Error iterasi/parsing ditelan diam-diam; hasil terpotong tanpa tanda.
  - **Perbaikan:** Panggil `rows.Err()` setelah loop dan tangani.

- [ ] **3.4 `GetPotensiBencana` 500 total jika satu baris memiliki geometri NULL**
  - **File:** `backend/internal/api/handlers/handlers.go:40-44`
  - **Masalah:** `ST_AsGeoJSON` bisa NULL; scan NULL ke `string` → seluruh endpoint error.
  - **Perbaikan:** Scan ke `*string`/`sql.NullString`; skip/fallback per baris.

- [ ] **3.5 Panic potensial `resp.Choices[0]`**
  - **File:** `backend/internal/ai/nlp.go:62,122`
  - **Masalah:** Jika model mengembalikan `choices` kosong → index out of range.
  - **Perbaikan:** Cek `len(resp.Choices) == 0` sebelum akses.

- [ ] **3.6 Cron jalan dobel saat server start + tanpa dedup berita**
  - **File:** `backend/internal/worker/cron.go:19-30`
  - **Masalah:** Jobs dieksekusi langsung di `main` (`go fetcher.Fetch...`) DAN terjadwal setiap 6 jam; jika server start tepat di jam kelipatan 6, data ditarik 2×. Insert berita tidak dedup → `kabar_kelud` membengkak 2 baris/siklus selamanya.
  - **Perbaikan:** Jalankan sekali lewat cron (mis. `AddFunc` dengan waktu start segera), tambahkan dedup (UNIQUE constraint/cek judul+sumber+interval).

- [ ] **3.7 JSON response AI bisa tidak valid secara konsisten**
  - **File:** `backend/internal/ai/nlp.go:121-127`
  - **Masalah:** `GenerateNewsSummary` mengandalkan output JSON persis; model kadang membungkus dengan ```json fence → parse gagal → berita hilang tanpa retry.
  - **Perbaikan:** Strip code fence, tambah retry/fallback sekali.

## 4. 🟠 SEDANG — Database & Migrasi

- [ ] **4.1 Migrasi menghapus data: `DROP TABLE potensi_bencana CASCADE`**
  - **File:** `backend/scripts/migrate.go:77`
  - **Masalah:** Menjalankan "migrasi" pada DB terisi = data hilang permanen (plus cascade). Script migrasi seharusnya idempotent & non-destruktif.
  - **Perbaikan:** Hapus `DROP TABLE`; pindahkan ke script seed/reset terpisah yang eksplisit meminta konfirmasi; evaluasi pemakaian tool migrasi (golang-migrate/atlas).

- [ ] **4.2 Dimensi vektor inkonsisten: 1024 (kode) vs 384 (dokumen)**
  - **File:** `backend/scripts/migrate.go:88` (1024), `Schema.md:22` (384), `README.md:75` (384)
  - **Masalah:** Runtime benar memakai 1024 (Cohere `embed-multilingual-v3.0`), tapi docs bilang 384 — kalau ada yang mengikuti docs, insert vektor gagal.
  - **Perbaikan:** Sinkronkan dokumen ke 1024, atau komentar di migrate.go menjelaskan dimensi & model embedding.

- [ ] **4.3 Tidak ada index untuk query spasial & vektor**
  - **File:** `backend/scripts/migrate.go:63-101` (schema), `ai_handler.go:97` (`ORDER BY embedding <=> $1`), `handlers.go` (SELECT semua geometri)
  - **Masalah:** `Schema.md:26-27` menjanjikan GiST pada geometri & HNSW pada embedding, tapi tidak dibuat di migrate.go → full scan 1024-dimensi per query.
  - **Perbaikan:** Tambah `CREATE INDEX ... USING GiST (geometri)` dan `... USING hnsw (embedding vector_cosine_ops)`.

- [ ] **4.4 `insert_hospitals.go` tidak mengisi `embedding`**
  - **File:** `backend/scripts/insert_hospitals.go:32-35`
  - **Masalah:** RSUD Gambiran & RSKK Pare punya `embedding NULL` → tidak pernah muncul di hasil semantic search.
  - **Perbaikan:** Generate embedding saat seed (atau backfill script).

- [ ] **4.5 `cleanup.go` memakai `MIN(id_potensi::text)` untuk "baris paling awal"**
  - **File:** `backend/scripts/cleanup.go:33-38`
  - **Masalah:** `MIN` pada representasi teks UUID = urutan lexicographic, bukan urutan insert; berisiko menghapus baris yang salah.
  - **Perbaikan:** Gunakan `created_at`/`ctid` atau kolom urutan insert.

- [ ] **4.6 Dockerfile DB rawan gagal build**
  - **File:** `backend/docker/Dockerfile-db`
  - **Masalah:** `postgresql-16-pgvector` (paket Debian) kemungkinan tidak tersedia di repo *bookworm* basis image `postgis/postgis:16-3.4` → build bisa gagal.
  - **Perbaikan:** Gunakan image `pgvector/pgvector:pg16` + install PostGIS, atau verifikasi ketersediaan paket & pin versi.

## 5. 🟠 SEDANG — Frontend

- [ ] **5.1 `KeludMapView.svelte` — layer tidak dibersihkan & refetch berulang**
  - **File:** `frontend/src/lib/components/KeludMapView.svelte:196-398`
  - **Masalah:** (a) `routeLayer` & `originMarker` tidak pernah dihapus saat toggle mati → layer menumpuk; (b) rute di-*refetch* setiap `$effect` berjalan; (c) koordinat awal/akhir diduplikasi di polyline (server sudah mengembalikan endpoint).
  - **Perbaikan:** Tambah cabang `else` yang menghapus layer; cache hasil fetch per pasangan koordinat; buang titik duplikat.

- [ ] **5.2 Tile Google Satellite tanpa API key**
  - **File:** `frontend/src/lib/components/KeludMapView.svelte:73-76`
  - **Masalah:** `mt1.google.com` dipakai tanpa key → melanggar ToS Google & rawan diblokir/berhenti bekerja.
  - **Perbaikan:** Ganti ke layanan tile berlisensi (MapTiler/OSM/HERE) dengan key via env.

- [ ] **5.3 Filter blocklist/allowlist query false-positive & hanya di client**
  - **File:** `frontend/src/routes/search/+page.svelte:14-23`, `frontend/src/routes/safe-routes/+page.svelte:175-181` (regex di-duplikasi di 2 file)
  - **Masalah:** Query sah seperti "...kurang jauh" kena blocklist (`tambah|kurang|dibagi|dikali`) → muncul pesan "Anomali" yang membingungkan warga. Regulasi ini juga mudah dilewati (client-side).
  - **Perbaikan:** Satu modul bersama, evaluasi ulang kata-kata, dan pindahkan validasi sebenarnya ke backend.

- [ ] **5.4 Build produksi tidak berfungsi: `adapter-auto` tanpa platform adapter**
  - **File:** `frontend/svelte.config.js:9`, `frontend/package.json` (`build` script)
  - **Masalah:** Tidak ada adapter (node/vercel/netlify) di dependencies → `npm run build` gagal kecuali terdeteksi platform serveless saat build.
  - **Perbaikan:** Tambah `@sveltejs/adapter-node` (bila deploy VM/Docker) atau adapter platform yang dituju.

- [ ] **5.5 Input koordinat rute tidak divalidasi**
  - **File:** `frontend/src/routes/api/route/+server.ts:11-17`
  - **Masalah:** `split(',')` tanpa cek jumlah elemen; `parseFloat` bisa menghasilkan `NaN` yang diteruskan ke ORS/OSRM.
  - **Perbaikan:** Validasi `Number.isFinite` untuk 4 nilai; balas 400 bila invalid.

## 6. 🟡 RENDAH — Best Practice & Kebersihan

- [ ] **6.1 Tidak ada pengujian sama sekali**
  - Tidak ada `_test.go` / `.spec.ts`. `Rules.md §4` bahkan meminta load testing untuk endpoint routing.
  - **Perbaikan:** Mulai dengan unit test handler + integrasi migrasi, lalu load test `/api/search` & `/api/potensi`.

- [ ] **6.2 Tidak ada CI / lint otomatis**
  - Tidak ada workflow CI; `go vet`, `svelte-check` (`npm run check`), eslint tidak dijalankan otomatis.
  - **Perbaikan:** Tambah GitHub Actions: `go vet ./...`, `go test ./...`, `npm run check`, build Docker.

- [ ] **6.3 Docker anti-pattern**
  - **File:** `frontend/Dockerfile`, `backend/docker/Dockerfile-backend`, `docker-compose.yaml`
  - **Masalah:** Container jalan sebagai root; tanpa `HEALTHCHECK`; frontend menjalankan **dev server** (`npm run dev`) sebagai "produksi"; `env_file: ./backend/.env` membuat compose gagal bila file tidak ada; `version: '3.8'` obsolete.
  - **Perbaikan:** Multi-stage build + non-root user + production server; tambah healthcheck & `depends_on: condition: service_healthy`; buat `.env.example`.

- [ ] **6.4 Dokumentasi tidak sinkron dengan implementasi**
  - **File:** `General.md:17,31`, `Architecture.md:15,19,54`, `Schema.md:41,84`, `Rules.md:9`, `README.md`
  - **Masalah:** Proposal mewajibkan port 6000/6001 & `pgRouting`, implementasi memakai 5180/5181 & routing eksternal ORS/OSRM. README menyebut DeepSeek/HuggingFace/NASA "tiap 1 jam" + `GET /api/alert`, padahal kode memakai env `AI_API_KEY`/`AI_BASE_URL`/`AI_MODEL`/`COHERE_API_KEY`, cron 6 jam, dan `/api/alert` **tidak pernah dipanggil frontend** (fitur mati).
  - **Perbaikan:** Update dokumen ke realita, atau hapus fitur yang tidak dipakai; buat `docs/` arsitektur aktual.

- [ ] **6.5 Konvensi env pecah & tidak ada `.env.example`**
  - **File:** `README.md:65-66` (menyuruh `DEEPSEEK_API_KEY`/`HF_TOKEN`), kode membaca key berbeda.
  - **Perbaikan:** Satu kontrak env yang disepakati + `.env.example` di root/backend/frontend.

- [ ] **6.6 Kualitas kode frontend**
  - **File:** banyak di `src/routes/*` & `src/lib/*`
  - **Masalah:** `any` di mana-mana; `alert()` untuk error UX (`safe-routes/+page.svelte:179,258,265`); `navigator.clipboard` tanpa handling rejection (`safe-routes:563`); inline SVG WhatsApp di `Footer.svelte` di-duplikasi 2×; `i18n.t()` mengembalikan raw key saat typo (silent fail) — beberapa key `info.filter.*` di-cast `as any` yang menutupi error; `project.inlang/settings.json` kosong (dependensi mati).
  - **Perbaikan:** Ketik data dengan interface/type; komponen SVG reusable; tipe aman untuk i18n; hapus dep mati.

- [ ] **6.7 Aksesibilitas**
  - **File:** modal di `+page.svelte`, `news/+page.svelte`, `safe-routes/+page.svelte`, `danger-map/+page.svelte`
  - **Masalah:** Modal tidak ada focus trap, `Escape` tidak ditutup, backdrop pakai `onclick` pada `<div>` (hanya di-suppress via `svelte-ignore` a11y).
  - **Perbaikan:** Gunakan dialog semantics / library (bits-ui sudah ada di dependencies) dengan focus trap & keyboard handler.

- [ ] **6.8 Logging tidak terstruktur & mixed up**
  - Backend memakai `log.Printf` tersebar tanpa level/konteks (no slog); error handler menutup detail di satu tempat tapi bocor di tempat lain (lihat 2.3).
  - **Perbaikan:** Migrasi ke `log/slog` dengan level, tambahkan request ID.

- [ ] **6.9 Skrip ad-hoc tanpa standar**
  - **File:** `backend/scripts/*` (6 skrip dengan `//go:build ignore`, duplikasi koneksi DB & load `.env` setiap file)
  - **Masalah:** Tanpa framework CLI/migrasi; mudah salah jalan & tidak konsisten (lihat 4.4, 4.5).
  - **Perbaikan:** Konsolidasi jadi satu command (mis. `cmd/seed` dengan subcommand) memakai koneksi/config bersama.

