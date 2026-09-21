# DOCUMENT CONTEXT FOR AI AGENT (ANTIGRAVITY IDE) - LOKARI

**Project Code Name:** LOKARI (Platform WebGIS Dinamis Berbasis Zero-Admin dan Spatial AI)
**Lokasi Fokus:** Desa Jarak, Kec. Plosoklaten, Kab. Kediri (Lereng Gunung Kelud).

---

## 1. Visi Utama Proyek Berdasarkan Proposal Akademik
LOKARI adalah solusi mitigasi bencana alam struktural/digital yang ditujukan untuk menjawab minimnya infrastruktur informasi desa yang sering terbengkalai karena perangkat desa *gaptek* atau sibuk. Solusinya adalah membangun WebGIS *Zero-Admin*—sistem yang bekerja mandiri tanpa perlu manusia memasukkan data.

## 2. Kebutuhan Sistem Inti
1. **WebGIS (SvelteKit + Leaflet):** Memetakan Kawasan Rawan Bencana (KRB) Erupsi Gunung Kelud, alur lahar dingin Kali Ngobo, posko pengungsian, dan fasilitas pendukung.
2. **Zero-Admin Backend (Golang):** Server secara berkala (melalui eksekusi latar belakang) menarik data dari:
   - PVMBG / Magma Indonesia (Status Gunung Api, Amplitudo Tremor).
   - BMKG (Intensitas Curah Hujan).
3. **Spatial AI - NLP Summarizer:** Data angka mentah di atas dikonversi otomatis menjadi satu kalimat peringatan darurat berbahasa Indonesia (ditampilkan sebagai *Alert Banner*).
4. **Spatial AI - Semantic Search & Routing:** Integrasi bahasa alami untuk navigasi. Menggunakan ekstensi **pgvector** di PostgreSQL untuk membedah intent pencarian rute, kemudian dilanjutkan dengan eksekusi algoritma Dijkstra melalui **pgRouting** (PostGIS) untuk menggambar garis terpendek ke posko aman, yang menghindari jalanan terputus akibat zona lahar.

## 3. Resolusi Konsep (Tim vs Proposal Akademik)
Dokumen `General.md` awal sempat menyebutkan tumpukan teknologi NoSQL (MongoDB) serta desain UI berupa *Hero Carousel* dengan video *autoplay*.

**Penyesuaian Wajib:**
- Demi meluluskan konversi 10 SKS mata kuliah tim, **Proposal Akademik adalah Hukum Tertinggi**.
- Basis Data MongoDB diganti menjadi **PostgreSQL + PostGIS** karena kebutuhan spesifik *pgvector* dan *pgRouting*.
- UI utama diubah dari panggung presentasi wisata menjadi **Dashboard Navigasi Keselamatan (Active and Simplified)**. Tata letak *frontend* yang telah diprototipe menggunakan Base44 (React) akan dikonversi (di-refactor) sintaksnya ke dalam bahasa **SvelteKit (.svelte)** tanpa menggunakan *library* pihak ketiga Base44.

## 4. Pembagian Tugas & Siklus Sprint (Agile Scrum)
Proyek berjalan dalam 4 Sprint (15-30 September 2026).
- **Product Owner:** Analis Sistem & Yayasan Sagasitas.
- **Scrum Master:** Bertanggung jawab atas pembagian tugas antar 15 anggota mahasiswa lintas prodi (Teknik Informatika, Sistem Informasi, Administrasi Negara).
- **Development Team (Fokus Anda: Server Infrastructure Engineer):** Bertanggung jawab memastikan peladen Linux sehat, Port Nginx (*Frontend* 6000, *Backend* 6001) tersambung baik, penanganan *PM2/Docker*, serta *stress testing* atas beban kerja Golang dan PostgreSQL.

---
*Status: DIBARUI DAN DISESUAIKAN DENGAN PROPOSAL AKADEMIK*
