# LOKARI (Lokasi Relokasi & Evakuasi Mandiri)

Selamat datang di repositori **LOKARI**! 
LOKARI adalah sebuah aplikasi pemetaan berbasis *website* yang cerdas dan interaktif, dirancang khusus untuk memandu dan memitigasi dampak bencana erupsi **Gunung Kelud**, khususnya bagi warga Desa Jarak, Kabupaten Kediri dan sekitarnya.

Aplikasi ini dilengkapi dengan fitur:
- **Peta Interaktif** dengan zona bahaya (KRB) langsung dari data satelit.
- **Pencarian Semantik AI (Vektor)** untuk mencari posko pengungsian atau fasilitas kesehatan terdekat dengan bahasa alami.
- **Rute Evakuasi Cerdas** yang secara otomatis menghindari Zona Merah.
- **Peringatan Dini Bencana** yang menarik data *real-time* dari satelit NASA EONET.

Proyek ini dibangun menggunakan **SvelteKit** (Frontend) dan **Go / Fiber** (Backend) dengan dukungan *database* **PostgreSQL (pgvector)**.

---

## Panduan Memulai (Untuk Pemula)

Jangan khawatir jika kamu baru pertama kali memegang *project* ini! Ikuti langkah-langkah di bawah ini secara berurutan.

### 1. Persiapan Alat Tempur
Sebelum memulai, pastikan kamu sudah meng-*install* aplikasi berikut di komputermu:
- [Git](https://git-scm.com/downloads) (Untuk mengunduh kodingan).
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (Wajib untuk menjalankan *database* dan lingkungan *server* secara instan).
- [Node.js](https://nodejs.org/) (Jika ingin menjalankan Frontend secara manual).
- [Go](https://go.dev/dl/) (Jika ingin menjalankan Backend secara manual).

### 2. Cara Mengunduh (Clone) Projek
Buka Terminal / Command Prompt / Git Bash, lalu jalankan perintah ini:
```bash
git clone https://github.com/USERNAME/LOKARI.git
cd LOKARI
```
*(Ganti tulisan `USERNAME` dengan nama akun GitHub tempat repositori ini berada).*

---

## Menjalankan Projek Menggunakan Docker (Cara Paling Mudah)

Dengan Docker, kamu tidak perlu repot *install* *database* secara manual. Cukup jalankan 1 perintah, semuanya akan menyala!

### Langkah-langkah:
1. Pastikan **Docker Desktop** sudah menyala (buka aplikasinya dan tunggu hingga ikon *engine* berwarna hijau).
2. Buka terminal di dalam folder `LOKARI`.
3. Jalankan perintah ajaib ini:
   ```bash
   docker compose up --build -d
   ```
4. Tunggu beberapa saat. Docker akan mengunduh dan merakit Frontend, Backend, dan Database secara otomatis.
5. Jika sudah selesai, buka *browser* dan ketik:
   - Frontend (Web LOKARI): `http://localhost:5180`
   - Backend API: `http://localhost:5181`

### Menyiapkan Tabel dan Mengisi Database (Migrasi & Seeder)
Jika kamu menjalankan proyek ini di laptop baru, *database* PostgreSQL di Docker masih sepenuhnya kosong. Kamu harus menjalankan perintah migrasi tabel terlebih dahulu sebelum mengisinya dengan data.

1. **Jalankan Migrasi Tabel Utama:**
   Membuat struktur tabel peta, titik lokasi, dan fitur AI.
   ```bash
   docker compose exec backend go run scripts/migrate.go
   ```
2. **Jalankan Migrasi Tabel Berita (Kabar Kelud):**
   Membuat tabel untuk menyimpan berita bencana.
   ```bash
   docker compose exec backend go run scripts/migrate_news.go
   ```
3. **Isi Database dengan Vektor AI (Seeder):**
   Memasukkan 58 titik kumpul dan posko dari file teks, lalu men-*generate* vektor AI-nya.
   ```bash
   docker compose exec backend go run scripts/seed_claude.go
   ```
Tunggu hingga proses ekstraksi AI selesai 100%!

> **Reset database (opsional):** untuk menghapus seluruh data lalu membangun ulang skema dan data dari nol, jalankan:
> ```bash
> docker compose exec backend go run scripts/reset.go
> ```
> Skrip ini **meminta konfirmasi `y/N`** sebelum menghapus apa pun, lalu menjalankan migrasi (`migrate.go` + `migrate_news.go`) dan menawarkan seed ulang.

> **Backfill embedding (opsional):** jika ada baris lama yang belum punya vektor (mis. hasil seed sebelum fitur AI embedding ada), lengkapi dengan:
> ```bash
> docker compose exec backend go run scripts/backfill_embedding.go
> ```
> Skrip ini hanya memproses baris dengan `embedding IS NULL` (idempotent).

---

## Menjalankan Projek Secara Manual (Tanpa Docker Penuh)

Jika kamu ingin ngoding dan melihat perubahannya secara langsung (*hot-reload*), gunakan cara ini:

### 1. Nyalakan Database Saja (Via Docker)
Kamu tetap butuh Docker HANYA untuk *database* PostgreSQL.
```bash
docker compose up -d db
```

### 2. Jalankan Backend (Go)
Buka terminal baru, masuk ke folder `backend`, lalu jalankan:
```bash
cd backend
go mod tidy
go run ./cmd/server/main.go
```
*(Backend akan menyala di port 5181).*

### 3. Jalankan Frontend (SvelteKit)
Buka terminal baru lagi, masuk ke folder `frontend`, lalu jalankan:
```bash
cd frontend
npm install
npm run dev
```
*(Frontend akan menyala di port 5173. Silakan buka http://localhost:5173).*

---

## 5. Deploy Frontend ke Vercel (Serverless)

> **Arsitektur deploy:** Hanya **frontend SvelteKit** yang serverless di Vercel.
> Backend Go (Fiber) tetap berjalan sebagai server biasa (VPS/Docker/Railway/Render/dsb.),
> dan frontend meneruskan `/api/*` ke backend lewat `src/hooks.server.ts`.

### 5.1 Prasyarat

1. **Backend Go sudah jalan di suatu host dengan URL publik HTTPS**
   (contoh: `https://lokari-api.example.com`). Pastikan endpoint
   `GET /api/health` bisa diakses publik.
2. **Database PostgreSQL + PostGIS + pgvector** yang bisa diakses backend.
   Cocok pakai **Neon** (supports PostGIS & pgvector) atau Supabase.
   Skema didaftarkan sekali dari lokal: `go run scripts/migrate.go` pada `backend/`.
3. **Akun Vercel** dan repo ini sudah di-push ke GitHub/GitLab.

### 5.2 Langkah Deploy (Dashboard Vercel)

1. **Import project** dari repo GitHub → pilih **Root Directory: `frontend`**.
   - Framework terdeteksi otomatis (SvelteKit) berkat `frontend/vercel.json`.
   - Build command otomatis `npm run build` (adapter-vercel menghasilkan output serverless).
2. Set **Environment Variables** (Settings → Environment Variables):
   | Variabel | Keterangan |
   |---|---|
   | `BACKEND_URL` | URL publik backend Go, mis. `https://lokari-api.example.com` |
   | `ORS_API_KEY` | API key OpenRouteService (untuk `/api/route`) |
3. **Deploy.** Selesai.

> ⚠️ **Keamanan:** API key OpenRouteService WAJIB disimpan sebagai env
> `ORS_API_KEY` (server-only, dibaca di `src/routes/api/route/+server.ts`),
> jangan pernah di-hardcode di source code.

### 5.3 Cara kerja `/api/*` di Vercel

- `/api/route` → ditangani endpoint SvelteKit sendiri (`src/routes/api/route/+server.ts`,
  routing ORS/OSRM).
- `/api/potensi`, `/api/news`, `/api/search`, `/api/alert`, `/api/health` →
  request jatuh ke fungsi SSR → `src/hooks.server.ts` meneruskan ke `BACKEND_URL`
  (server-side, jadi tidak ada masalah CORS).
- Var env `BACKEND_URL` dibaca runtime via `$env/dynamic/private`
  (default fallback `http://localhost:5181` hanya untuk dev lokal).

### 5.4 Node.js runtime

- `frontend/vercel.json` + `frontend/svelte.config.js` mem-pin runtime
  **`nodejs22.x`** (Vercel mendukung 20/22/24). Jangan pindah ke `edge` —
  `hooks.server.ts` dan `+server.ts` memerlukan Node runtime.

### 5.5 Deploy lewat CLI (opsional)

```bash
cd frontend
npx vercel link          # hubungkan ke project Vercel (root directory: frontend)
npx vercel env add BACKEND_URL   # tambahkan variabel env
npx vercel --prod
```

### 5.6 GitHub Actions (CI/CD)

Workflow di `.github/workflows/` dibatasi hanya berjalan di repo tim
(`lokari-unesa/lokari`) lewat `if: github.repository == ...` — repo fork tidak
menjalankannya. Secrets yang dibutuhkan di repo tersebut: `VERCEL_TOKEN`,
`ORG_ID`, `PROJECT_ID`.

---

## 6. Panduan Deployment (Hosting ke VPS Server)

Jika aplikasi LOKARI sudah siap *go-public* dan ingin di- *hosting* ke server VPS (misal: AWS, DigitalOcean, atau IDCloudHost), ada 2 metode yang bisa kamu pilih:

### Opsi A: Deployment Penuh via Docker (Direkomendasikan)
Cara ini paling bersih dan persis sama seperti saat kita menjalankannya di laptop lokal.

1. Beli VPS berbasis Linux (Ubuntu).
2. *Install* Git dan Docker di VPS.
3. *Clone* repositori ini ke VPS.
4. Ubah sedikit *file* `.env` jika ada URL API produksi yang perlu diubah.
5. Jalankan perintah andalan:
   ```bash
   docker compose up --build -d
   ```
6. Opsional: *Install* NGINX di VPS (luar Docker) sebagai *Reverse Proxy* untuk mengarahkan Domain (misal: `lokari.com`) ke `localhost:5180` dan memasang SSL (HTTPS) menggunakan *Certbot/Let's Encrypt*.

### Opsi B: NGINX + PM2 (Frontend/Backend) & Docker (Database Saja)
Cara tradisional ini cocok jika kamu ingin memantau *service* menggunakan PM2.

1. **Database:**
   Di VPS, jalankan HANYA *database* menggunakan Docker.
   ```bash
   docker compose up -d db
   ```
2. **Backend (Go):**
   - Lakukan *build* aplikasi Go: `cd backend && go build -o lokari-app ./cmd/server/main.go`
   - *Install* PM2 (via Node.js).
   - Jalankan backend dengan PM2: `pm2 start ./lokari-app --name "lokari-backend"`
3. **Frontend (SvelteKit Node Adapter):**
   - Lakukan *build* Svelte: `cd frontend && npm install && npm run build`
   - Ini akan menghasilkan folder `build/`.
   - Jalankan *frontend* dengan PM2: `pm2 start build/index.js --name "lokari-frontend"`
4. **NGINX Reverse Proxy:**
   - Konfigurasi `/etc/nginx/sites-available/lokari` untuk meneruskan lalu lintas:
     - `lokari.com` → `http://localhost:5180` (port SvelteKit)
     - `api.lokari.com` → `http://localhost:5181` (port Go Backend)
   - Pasang sertifikat SSL.

---
*Dibuat oleh Tim Universitas Negeri Surabaya untuk warga Desa Jarak.*
