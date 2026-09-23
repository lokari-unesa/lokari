# LOKARI (Lokasi Aman & Rute Evakuasi) 🌋

LOKARI adalah platform cerdas berbasis mitigasi bencana (khususnya Gunung Kelud untuk Desa Jarak, Kediri). Sistem ini memadukan **Database Spasial (PostGIS)**, **Zero-Admin Cron Worker**, dan **Spatial AI (Semantic Search & NLP Summarizer)** menggunakan `pgvector`, DeepSeek, dan Hugging Face.

Proyek ini dibangun menggunakan **SvelteKit (Frontend)** dan **Golang Fiber (Backend)**.

---

## Prasyarat (Prerequisites)
Sebelum menjalankan proyek ini, pastikan komputer/server Anda sudah terinstal:
1. **Git**
2. **Node.js** (Minimal versi 18.x) - *Untuk Frontend SvelteKit*
3. **Go (Golang)** (Minimal versi 1.21) - *Untuk Backend API*

---

## Langkah Instalasi

### 1. Clone Repository
```bash
git clone <URL_REPO_GITHUB_ANDA>
cd LOKARI
```

### 2. Setup Database (PostgreSQL + PostGIS + pgvector)
Anda bisa memilih salah satu dari dua cara di bawah ini untuk menghidupkan database:

#### Opsi A: Menggunakan Docker (Sangat Direkomendasikan)
Cara paling mudah tanpa perlu menginstal PostgreSQL secara manual. Pastikan Docker Desktop menyala.
1. Dari root proyek, jalankan database saja: `docker compose up -d db`
   *(atau `docker compose up -d --build` untuk full-stack dev, lihat bagian Docker di bawah)*
*(Database akan berjalan di port host `5433`).*

#### Opsi B: Instalasi Manual (Native / Tanpa Docker)
Jika Anda tidak menggunakan Docker, instal secara manual sesuai Sistem Operasi Anda:

**Untuk Windows:**
1. Unduh dan instal PostgreSQL via *EnterpriseDB installer*.
2. Setelah instalasi selesai, buka aplikasi bawaannya yaitu **Application Stack Builder**.
3. Pilih server PostgreSQL Anda pada menu *dropdown*, lalu klik Next.
4. Buka kategori **Spatial Extensions**, centang **PostGIS** dan ikuti proses instalasinya sampai selesai.
*(Catatan: Menginstal ekstensi `pgvector` secara manual di Windows cukup rumit karena mewajibkan kompilasi C++ Visual Studio. Untuk pengembangan di Windows, sangat disarankan beralih memakai Docker).*

**Untuk Server Linux (Ubuntu/Debian):**
Buka terminal server Anda dan jalankan deretan perintah berikut (ganti angka `16` sesuai versi PostgreSQL yang Anda inginkan):
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo apt install postgis postgresql-16-postgis-3
sudo apt install postgresql-16-pgvector
```

---

### 3. Konfigurasi Backend (Golang)
Masuk ke folder `backend`, salin contoh env lalu isi nilai aslinya:
```bash
cp backend/.env.example backend/.env   # lalu edit nilainya
```

Isi `backend/.env` (file ini TIDAK di-commit):
```env
PORT=5181
# Jika pakai Docker (host port 5433 -> container 5432):
DATABASE_URL=postgresql://postgres:root@127.0.0.1:5433/lokari_db?sslmode=disable
# Jika install manual lokal (Port 5432):
# DATABASE_URL=postgresql://postgres:root@127.0.0.1:5432/lokari_db?sslmode=disable

# AI untuk rangkuman berita/alert (provider OpenAI-compatible: OpenRouter/Groq/DeepSeek)
AI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
AI_BASE_URL=https://openrouter.ai/api/v1
AI_MODEL=openrouter/free

# AI untuk vector embeddings (semantic search, dimensi 1024)
COHERE_API_KEY=cohere_xxxxxxxxxxxxxxxxxxxxxxxx
```

**Frontend (opsional):** tidak wajib punya `.env` — default sudah `VITE_PORT=5180` dan proxy `/api` → `http://localhost:5181`. Salin `frontend/.env.example` ke `frontend/.env` hanya bila ingin mengubahnya.

**Docker Compose — dua mode (disarankan):**

| Mode | Perintah | Keterangan |
|---|---|---|
| **Development** (default) | `docker compose up -d --build` | Hot reload: backend `air`, frontend `vite dev` + bind mount kode |
| **Production / Build** | `docker compose -f docker-compose.yaml up -d --build` | Backend: binary statis (alpine). Frontend: bundle `@sveltejs/adapter-node` (`node build`) |

> Mode development dimuat otomatis lewat `docker-compose.override.yaml`. Untuk production murni, jalankan tanpa override dengan `-f docker-compose.yaml`.
> Interpolasi port & kredensial DB bisa diatur lewat `.env` di root (opsional, contoh: `.env.example`). Rahasia backend tetap di `backend/.env`.
> Di mode production tidak ada vite proxy, jadi `/api/*` diteruskan ke backend Go via `src/hooks.server.ts`.

**Instal dependensi Go:**
```bash
go mod tidy
```

**Inisialisasi Database (Hanya dilakukan 1x di awal):**
Ini akan secara otomatis membuat tabel, mengaktifkan `postgis`, dan membuat kolom `vector(1024)`.
```bash
go run scripts/migrate.go
```

**Jalankan Backend — mode Development (hot-reload otomatis dengan air):**
```bash
air
```
*(Pastikan `air` terinstal: `go install github.com/air-verse/air@latest`).*

**Jalankan Backend — mode Production (binary):**
```bash
go build -o main ./cmd/server/main.go && ./main
```
*(Backend akan berjalan di `http://localhost:5181`. Zero-Admin Cron Job akan otomatis menarik data dari NASA setiap 1 jam).*

---

### 4. Konfigurasi Frontend (SvelteKit)
Buka tab terminal/PowerShell **baru**, lalu arahkan ke folder frontend:
```bash
cd ../frontend
```

**Instal dependensi Node.js:**
```bash
npm install
```

**Jalankan Frontend — mode Development (vite dev server):**
```bash
npm run dev
```
*(Frontend akan berjalan di `http://localhost:5180`. Tekan tombol `o` di terminal, atau klik tautan tersebut untuk membukanya di browser).*

**Jalankan Frontend — mode Production (bundle adapter-node):**
```bash
npm run build && npm run preview
```
*(Preview menjalankan hasil build; `/api/*` diteruskan ke backend Go via `src/hooks.server.ts` — backend lokal harus jalan di `http://localhost:5181`. Nilai target bisa diganti lewat env `BACKEND_URL`).*

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

## Catatan Khusus Modul AI
1. **GET `/api/alert`**: Rute ini akan menarik data satelit dari **NASA EONET**, kemudian diringkas menggunakan **DeepSeek API** menjadi bahasa Indonesia yang ramah warga.
2. **POST `/api/search`**: Menerima input kalimat acak (contoh: *"posko evakuasi lahar terdekat"*), mengubahnya jadi vektor memakai **Hugging Face**, lalu dicarikan dengan jarak semantik `<=>` dari ekstensi `pgvector`.

✨ **Semangat untuk tim LOKARI Universitas Negeri Surabaya (UNESA)!** ✨
