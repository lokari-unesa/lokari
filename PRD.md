# Product Requirements Document (PRD) - LOKARI

> Dokumen ini diselaraskan dengan **fitur yang benar-benar terimplementasi**
> (Oktober 2026). Setiap requirement berlabel status terhadap proposal awal:
> ✅ terpenuhi · ⚠️ diubah (dengan alasan) · ❌ diganti/dibatalkan.

## 1. Ringkasan Eksekutif

LOKARI adalah platform *Web-based Geographic Information System* (WebGIS)
dinamis dengan pendekatan *Zero-Admin* dan *Spatial AI* untuk mitigasi
bencana erupsi Gunung Kelud dan aliran lahar dingin Kali Ngobo di Desa
Jarak, Kecamatan Plosoklaten, Kabupaten Kediri.

### 1.1 Visi

Menyediakan infrastruktur digital keselamatan yang aktif, akurat, dan dapat
dijangkau seluruh lapisan masyarakat pedesaan untuk meminimalkan risiko
korban jiwa saat bencana erupsi maupun aliran lahar terjadi.

### 1.2 Misi

"Menghubungkan informasi dengan ruang, memudahkan akses menuju lokasi."

## 2. Target Pengguna & Evaluasi Penerimaan (TAM/UTAUT)

Platform dirancang untuk warga yang awam teknologi. Keberhasilan diukur
dengan *Technology Acceptance Model* (TAM):

- **Persepsi Kemudahan:** antarmuka aktif dan sederhana, bebas kerumitan
  operasional. Ditambah dukungan dwibahasa ID/EN dan mode terang/gelap.
- **Persepsi Kebermanfaatan:** keakuratan rute evakuasi, kecepatan
  peringatan dini (notifikasi push browser), efisiensi waktu mencari posko.

## 3. Fitur Fungsional

### 3.1 Zero-Admin Otomatis — ✅ terpenuhi (sumber data disesuaikan)

- **REQ-1.1:** *Backend* (Golang) menjalankan *background worker* dual-loop:
  **hot loop tiap 30 detik** (gempa BMKG terbaru + status Gunung Kelud dari
  PVMBG MAGMA) dan **cold loop tiap 6 jam** (laporan harian MAGMA + event
  vulkanik NASA EONET). ⚠️ *Deviasi dari proposal*: proposal menyebut API
  PVMBG (status + tremor) dan BMKG (curah hujan) tiap 10 menit. Realita:
  tidak ada endpoint API resmi PVMBG yang stabil sehingga status Kelud
  diambil via scrape halaman MAGMA, dan BMKG yang dipakai adalah data
  **gempa** (bukan curah hujan) dengan filter magnitudo ≥ 3.5 + radius
  dinamis dari Kelud.
- **REQ-1.2:** Perangkat desa tidak perlu input manual; status bahaya,
  berita, dan notifikasi diproses otomatis di server. ✅

### 3.2 Pemetaan WebGIS — ✅ terpenuhi sebagian, ⚠️ sumber geometri disederhanakan

- **REQ-2.1:** ⚠️ Visualisasi zona bahaya KRB 1/2/3 sebagai **lingkaran
  Leaflet hardcoded** (radius 5/10/15 km dari kawah), bukan poligon KRB dari
  database (tabel `kawasan_rawan` tidak pernah dibuat — poligon KRB resmi
  tidak tersedia sebagai data geospasial). Warna aktual: `#EA580C` (bahaya),
  `#16A34A` (posko), `#0284C7` (kesehatan), `#92400E` (lahar).
- **REQ-2.2:** ⚠️ Alur lahar Kali Ngobo sebagai polyline/garis simulasi
  hardcoded di frontend, bukan `LINESTRING` dari database.
- **REQ-2.3:** ✅ Titik (*POINT*) posko pengungsian, titik kumpul, dan
  fasilitas kesehatan dari `potensi_bencana` via `GET /api/potensi`.

### 3.3 AI: Ringkasan Berita Otomatis — ⚠️ dipindah dari Alert Banner ke berita

- **REQ-3.1:** ✅ Mesin AI (provider OpenAI-compatible via
  `AI_API_KEY`/`AI_BASE_URL`/`AI_MODEL`) mengonversi data mentah
  BMKG/MAGMA/EONET menjadi ringkasan berita Bahasa Indonesia, disimpan di
  `kabar_kelud` (dedupe `ON CONFLICT (sumber, judul)`). Ada teks fallback
  deterministik bila AI gagal.
- **REQ-3.2:** ❌ *Alert Banner* real-time **tidak diimplementasikan**.
  Endpoint `GET /api/alert` ada di backend tetapi tidak pernah dipanggil
  frontend (fitur mati). Fungsinya digantikan **kartu status Gunung Kelud**
  di beranda yang membaca `GET /api/kelud/status` (sumber kebenaran:
  `monitor_state`), plus **notifikasi Web Push** otomatis untuk setiap
  berita penting.

### 3.4 AI: Semantic Search & Rute Cerdas — ✅ search terpenuhi, ❌ routing diganti

- **REQ-4.1:** ✅ *Semantic search* memakai PostgreSQL + **pgvector**:
  embedding Cohere `embed-multilingual-v3.0` (1024 dimensi, index HNSW).
  Backend me-retrieve Top 15 (`POST /api/search`), frontend me-rerank
  berdasar jarak fisik.
- **REQ-4.2:** ✅ Pertanyaan bahasa sehari-hari ("lari ke mana kalau ada
  lahar?") dikonversi ke vektor dan dicocokkan dengan lokasi terdekat.
- **REQ-4.3:** ❌ **pgRouting/Dijkstra tidak diimplementasikan** (tidak
  pernah ada digitasi topologi jalan `jaringan_jalan` untuk Desa Jarak).
  Digantikan `GET /api/route` (SvelteKit): OpenRouteService dengan
  `avoid_polygons` + fallback OSRM, dengan 3 state UI (aman/fallback/bahaya).

### 3.5 Antarmuka Publik — ✅ terpenuhi + 5 halaman tambahan

- **REQ-5.1:** ✅ Peta interaktif SvelteKit + Leaflet (bukan full-screen —
  lihat `Design.md`).
- **REQ-5.2:** ✅ *Pop-up* interaktif: nama, kapasitas, kategori,
  rekomendasi faskes terdekat, dan tombol rute evakuasi.
- **Tambahan di luar proposal**: halaman `/news` (Kabar Kelud),
  `/search`, `/safe-routes`, `/danger-map`, `/about`; i18n ID/EN; mode
  gelap; Web Push notification.

## 4. Kebutuhan Non-Fungsional

- **NFR-1 (Kinerja):** pencarian semantik dibatasi Top 15 + index HNSW;
  routing didelegasikan ke layanan eksternal agar DB tidak terbebani.
- **NFR-2 (Keandalan):** semua HTTP keluar ber-timeout; worker dan handler
  bertahan tanpa DB (`nil` pool diguard); error internal dilog di server,
  client menerima pesan generik; endpoint AI ber-rate-limit 30 req/menit/IP.

## 5. Realisasi Fase Pengerjaan

Pengembangan berjalan iteratif (Sept–Okt 2026) di luar skema 4 sprint
proposal: fondasi DB + backend Zero-Admin, lalu frontend + semantic search,
diikuti pengerasan (rate limit, dedupe berita, cleanup subscription push,
dual-loop cron, monitor status Kelud) dan deploy (Vercel + Docker).

---

*Status: SELARAS DENGAN IMPLEMENTASI (Okt 2026). Detail teknis: `docs/`.*
