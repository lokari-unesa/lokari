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
1. Masuk ke folder backend: `cd backend`
2. Jalankan mesin database: `docker-compose up -d`
*(Database akan berjalan di port `5433`)*
*Perintah stop/start: `docker stop lokari_db_vector` / `docker start lokari_db_vector`.*

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
Masuk ke folder `backend`, buat file bernama `.env` (atau salin jika sudah ada), dan isi dengan:
```env
PORT=5181
# Jika pakai Docker (Port 5433):
DATABASE_URL=postgresql://postgres:root@127.0.0.1:5433/lokari_db?sslmode=disable
# Jika install manual lokal (Port 5432):
# DATABASE_URL=postgresql://postgres:root@127.0.0.1:5432/lokari_db?sslmode=disable

DEEPSEEK_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
HF_TOKEN= # (Opsional jika API HuggingFace Anda belum terkena limit)
```

**Instal dependensi Go:**
```bash
go mod tidy
```

**Inisialisasi Database (Hanya dilakukan 1x di awal):**
Ini akan secara otomatis membuat tabel, mengaktifkan `postgis`, dan membuat kolom `vector(384)`.
```bash
go run scripts/migrate.go
```

**Jalankan Peladen (Server) Backend:**
```bash
go run cmd/server/main.go
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

**Jalankan Server Frontend:**
```bash
npm run dev
```
*(Frontend akan berjalan di `http://localhost:5180`. Tekan tombol `o` di terminal, atau klik tautan tersebut untuk membukanya di browser).*

---

## Catatan Khusus Modul AI
1. **GET `/api/alert`**: Rute ini akan menarik data satelit dari **NASA EONET**, kemudian diringkas menggunakan **DeepSeek API** menjadi bahasa Indonesia yang ramah warga.
2. **POST `/api/search`**: Menerima input kalimat acak (contoh: *"posko evakuasi lahar terdekat"*), mengubahnya jadi vektor memakai **Hugging Face**, lalu dicarikan dengan jarak semantik `<=>` dari ekstensi `pgvector`.

✨ **Semangat untuk tim LOKARI Universitas Negeri Surabaya (UNESA)!** ✨
