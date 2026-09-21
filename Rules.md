# Project Rules & Guidelines - LOKARI

## 1. Aturan Emas: Arsitektur Proposal Wajib Dipatuhi
Sebagai *AI Agent*, Anda **DILARANG KERAS** menggunakan tumpukan teknologi di luar proposal akademik. Kegagalan mematuhi ini berakibat fatal pada kelulusan proyek mahasiswa.
- **[DILARANG]** Menggunakan React, Next.js, Vue, atau tumpukan Base44 SDK.
- **[DILARANG]** Menggunakan MongoDB atau NoSQL lainnya.
- **[WAJIB]** Menggunakan SvelteKit (Frontend).
- **[WAJIB]** Menggunakan Golang (Backend).
- **[WAJIB]** Menggunakan PostgreSQL, PostGIS, pgvector, dan pgRouting (Database & Spasial).

## 2. Aturan Backend Golang (Zero-Admin & AI)
1. **Background Workers:** Modul penarik API dari BMKG & PVMBG harus berjalan independen menggunakan *Goroutines*. Jangan menunda respons API ke *client* (warga) gara-gara *server* sedang menarik data eksternal.
2. **Koneksi Database:** Pastikan modul Go menggunakan *driver* resmi PostgreSQL yang mendukung tipe data geometri/PostGIS (misalnya `jackc/pgx`).
3. **Pengolahan Teks (NLP):** Logika peringatan dini harus sederhana, stabil, dan bisa menangani skenario jika *model machine learning* sedang kelebihan beban (fallback).
4. **Security & Configuration:** Sama seperti aturan sebelumnya, *port* (e.g., `6000`, `6001`), DSN PostgreSQL, dan kunci API eksternal **WAJIB** disimpan dalam *environment variables* (`.env`). Tidak boleh ada *hardcode*.

## 3. Aturan Frontend SvelteKit
1. **Pemetaan Dasar:** Gunakan Leaflet.js atau Mapbox. Tangani konversi GeoJSON yang dikirim dari Golang agar dirender secara efisien di komponen Svelte.
2. **Migrasi dari .jsx ke .svelte:** Kita diizinkan mengonversi tata letak statis yang di-*generate* Base44 (React) ke SvelteKit selama ia mendukung konsep antarmuka *Active and Simplified*. Namun, semua kaitan dengan *library* khusus `@base44/sdk` wajib dibuang!
3. **Pemuatan Cepat (Mitigasi Kendala Internet):** Sesuai Subbab 3.3.2 proposal, ukuran beban *frontend* harus ditekan (optimasi gambar WebP, pengiriman data teringan) agar warga di pedesaan dengan sinyal lemah tetap dapat memuat peta interaktif secara responsif.

## 4. Manajemen Infrastruktur (Sesuai Peran Cloud Engineer)
- **Deployment:** Proses peluncuran kode akan menggunakan Nginx sebagai *reverse proxy*. *Frontend* dan *Backend* akan dirawat siklus hidupnya menggunakan Docker atau *Process Manager* seperti PM2.
- **Testing:** Selalu utamakan pengujian keandalan respon (API Load Testing) untuk *endpoint pgRouting* karena merupakan komputasi *database* yang berat.

---
*Status: SELARAS DENGAN PROPOSAL AKADEMIK*
