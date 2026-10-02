# UI/UX & Design System - LOKARI

> Menggambarkan antarmuka **sebagaimana diimplementasikan** (Oktober 2026).
> Pengguna utama tetap warga Desa Jarak yang awam teknologi, sehingga prinsip
> "Active and Simplified" dipertahankan.

## 1. Pendekatan Desain: "Active and Simplified UI/UX"

Sistem proaktif menampilkan data darurat otomatis (status Kelud, berita,
notifikasi push) dengan tampilan sesederhana mungkin: maksimal 1–2 klik
menuju info posko atau rute, simbol peta intuitif, dan teks Bahasa Indonesia
sehari-hari (dengan opsi English).

## 2. Kriteria Usability (TAM/UTAUT)

1. **Learnability:** geser peta dan klik ikon tanpa petunjuk; pencarian
   menerima bahasa sehari-hari.
2. **Efficiency:** info posko/rute ≤ 2 klik dari beranda atau peta bahaya.
3. **Effectiveness:** peta menjawab kebutuhan darurat seketika (zona bahaya
   + posko + rute dalam satu kanvas).
4. **Clarity/Visibility** (warna aktual di kode):
   - 🟠 `#EA580C` — zona bahaya.
   - 🟢 `#16A34A` — posko pengungsian / tempat aman.
   - 🔵 `#0284C7` — fasilitas kesehatan.
   - 🟤 `#92400E` — jalur lahar.
5. **Responsiveness:** SvelteKit responsif untuk ponsel; aset gambar WebP.
6. **Kejelasan simbol:** marker warna + emoji per kategori (🏥 kesehatan,
   🕌 ibadah, 🎓 pendidikan, 🏛️ balai/gedung).

Tambahan di luar proposal: **i18n ID/EN** (persist `lokari-locale`) dan
**mode terang/gelap** (persist `lokari-theme`).

## 3. Struktur Halaman (6 halaman, bukan peta full-screen)

Proposal memandatkan peta full-screen sebagai landing page. Realita: beranda
adalah **dashboard keselamatan** — hero (judul + 2 tombol aksi) di atas,
kartu status Kelud, aksi cepat, cuplikan peta, dan berita. Peta penuh hanya
di halaman khusus.

### 3.1 Beranda `/`

Hero + **kartu status Gunung Kelud** (sumber kebenaran tunggal:
`GET /api/kelud/status` dari `monitor_state`, bukan tebakan dari berita):

| Level | Label | Warna kartu |
|---|---|---|
| `Level I (Normal)` | NORMAL | hijau (`safe`) |
| `Level II (Waspada)` | WASPADA | kuning (`warning`) |
| `Level III (Siaga)` | SIAGA | biru (`secondary`) |
| `Level IV (Awas)` | AWAS | merah (`destructive`) |

Belum ada state → status netral "MEMANTAU". Teks imbauan per level
dwibahasa di store i18n (`status.normal` … `status.awas`).

### 3.2 Peta Bahaya `/danger-map`

Kanvas Leaflet + panel layer (bahaya/posko/kesehatan/lahar) + legenda.
Klik marker → popover: nama, dusun, deskripsi, kapasitas, kategori,
**rekomendasi faskes terdekat otomatis**, dan tombol rute ke
`/safe-routes` (koordinat via query string).

### 3.3 Jalur Evakuasi `/safe-routes`

Panel asal (GPS/manual) + tujuan (dari pencarian AI atau link peta).
Hasil rute memakai **3 state dari server** (`isSafe` + `status`):
aman (ORS + `avoid_polygons`), fallback (OSRM tanpa penghindaran),
bahaya. Tanpa ETA dari routing engine, ETA tampil "—" (bukan angka
karangan). Hasil di-cache per pasangan koordinat agar tidak refetch.

### 3.4 Kabar Kelud `/news` + Pencarian `/search`

- `/news`: 20 berita terbaru (`GET /api/news`), filter kategori
  (volcano/lahar/evac/weather/warning — slug EN dari DB, frontend
  memetakan label ID + slug EN), timestamp relatif (zona WIB).
- `/search`: pencarian semantik (`POST /api/search`, Top 15) + rerank
  jarak di client; hasil teratas tertaut ke `/safe-routes`.

### 3.5 Tentang `/about`

Profil tim, kemitraan (Sagasitas, Desa Jarak), dan sumber data resmi
(PVMBG MAGMA, BMKG, NASA EONET). Resolusi konflik lama: tidak ada Hero
Carousel video — edukasi menyatu di halaman ini dan beranda.

## 4. Isu terbuka

Tile peta Google Satellite dipakai **tanpa API key** (melanggar ToS,
rawan diblokir) — terlacak di TODO 5.2, belum diperbaiki.

---

*Status: SELARAS DENGAN IMPLEMENTASI (Okt 2026).*
