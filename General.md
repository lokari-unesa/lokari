# KONTEKS PROYEK - LOKARI

**Nama Proyek:** LOKARI (Lokasi Relokasi & Evakuasi Mandiri — Platform
WebGIS Dinamis Berbasis Zero-Admin dan Spatial AI)
**Lokasi Fokus:** Desa Jarak, Kec. Plosoklaten, Kab. Kediri (lereng Gunung Kelud).

---

## 1. Visi Proyek

LOKARI adalah solusi mitigasi bencana digital yang menjawab minimnya
infrastruktur informasi desa yang sering terbengkalai karena perangkat desa
sibuk atau awam teknologi. Jawabannya adalah WebGIS *Zero-Admin* — sistem
yang bekerja mandiri tanpa manusia memasukkan data: worker server menarik
data BMKG/MAGMA/EONET, AI meringkasnya jadi berita, dan warga menerima
notifikasi push otomatis.

## 2. Sistem Inti (sebagaimana diimplementasikan)

1. **WebGIS (SvelteKit 2 + Leaflet 1.9):** 6 halaman — beranda (kartu status
   Kelud + peta + berita), peta bahaya, jalur evakuasi, kabar Kelud,
   pencarian AI, tentang. Zona bahaya = lingkaran hardcoded (5/10/15 km),
   batas desa dari `desa_jarak.json`.
2. **Zero-Admin Backend (Go Fiber):** cron dual-loop — hot tiap 30 detik
   (gempa BMKG + status Kelud via scrape MAGMA), cold tiap 6 jam (laporan
   harian MAGMA + event NASA EONET). State di `monitor_state`, berita di
   `kabar_kelud`.
3. **Spatial AI:** ringkasan berita otomatis (provider OpenAI-compatible) +
   semantic search pgvector/Cohere 1024-d + rerank jarak di frontend.
4. **Routing:** `GET /api/route` (SvelteKit) — OpenRouteService dengan
   `avoid_polygons` + fallback OSRM. Bukan pgRouting.
5. **Notifikasi:** Web Push (VAPID) untuk setiap berita penting.

Detail: `docs/Architecture.md`, `docs/Schema.md`, `docs/API.md`.

## 3. Catatan sejarah: deviasi dari proposal akademik

Dokumen proposal awal (Sept 2026) mewajibkan stack PostgreSQL + PostGIS +
**pgRouting** + poligon KRB + port 6000/6001 + Nginx/PM2. Selama pengerjaan,
tiga hal itu tidak terwujud dengan alasan teknis yang terdokumentasi:

- **pgRouting → ORS/OSRM:** tidak pernah ada digitasi topologi jaringan
  jalan (`jaringan_jalan`) untuk Desa Jarak yang layak untuk Dijkstra +
  *dynamic hazard avoidance* di DB. Routing didelegasikan ke layanan
  eksternal; kolom `is_safe`/`cost` tidak pernah ada.
- **Poligon KRB → lingkaran hardcoded:** poligon KRB resmi tidak tersedia
  sebagai data geospasial; tabel `kawasan_rawan` tidak pernah dibuat.
- **Port 6000/6001 → 5180/5181** (frontend Docker/backend; 5173 untuk vite
  manual), **Nginx/PM2 →** proxy `hooks.server.ts` + deploy Vercel/Docker.

Aturan main sekarang ada di `Rules.md` — batasan stack yang masih berlaku
(SvelteKit, Golang, PostgreSQL; larangan React/MongoDB/Base44) tetap
dipertahankan, tetapi klaim yang terbukti tidak diimplementasikan tidak
lagi dinyatakan sebagai kewajiban.

## 4. Cara kerja tim

Pengembangan iteratif (Sept–Okt 2026) oleh tim mahasiswa lintas prodi
(Teknik Informatika, Sistem Informasi, Administrasi Negara) bersama Product
Owner (Analis Sistem & Yayasan Sagasitas). Keputusan arsitektur dicatat di
`TODO.md` (item selesai + isu terbuka) — baca itu sebelum mengubah alur
worker, skema DB, atau routing `/api/*`.

---

*Status: SELARAS DENGAN IMPLEMENTASI (Okt 2026).*
