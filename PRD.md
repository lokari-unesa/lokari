# Product Requirements Document (PRD) - LOKARI

## 1. Ringkasan Eksekutif (Executive Summary)
LOKARI adalah platform *Web-based Geographic Information System* (WebGIS) dinamis yang dirancang dengan pendekatan *Zero-Admin* dan *Spatial AI* untuk mitigasi bencana erupsi Gunung Kelud dan aliran lahar dingin Kali Ngobo di Desa Jarak, Kecamatan Plosoklaten, Kabupaten Kediri.

### 1.1 Visi Proyek
Menyediakan infrastruktur digital keselamatan yang aktif, akurat, dan dapat dijangkau oleh seluruh lapisan masyarakat pedesaan untuk meminimalkan risiko korban jiwa saat bencana erupsi maupun aliran lahar terjadi.

### 1.2 Misi Proyek
"Menghubungkan informasi dengan ruang, memudahkan akses menuju lokasi." (Sesuai kesepakatan tim dan selaras dengan misi mitigasi bencana non-struktural).

## 2. Target Pengguna & Evaluasi Penerimaan (TAM/UTAUT)
Platform dirancang khusus untuk warga yang awam teknologi (*gaptek*). Keberhasilan proyek diukur menggunakan kerangka *Technology Acceptance Model* (TAM):
- **Persepsi Kemudahan (Perceived Ease of Use):** Antarmuka *Active and Simplified*, bebas dari kerumitan operasional.
- **Persepsi Kebermanfaatan (Perceived Usefulness):** Keakuratan rute evakuasi, kecepatan informasi peringatan dini, dan efisiensi waktu warga dalam mencari posko pengungsian.

## 3. Fitur Fungsional Utama (Sesuai Proposal)

### 3.1 Pendekatan Zero-Admin Otomatis
- **REQ-1.1:** Sistem *backend* (Golang) harus memiliki *background worker* yang menarik data parameter kebencanaan dan cuaca secara *real-time* dari API resmi PVMBG dan BMKG.
- **REQ-1.2:** Perangkat desa tidak perlu melakukan *input* data manual di *admin panel* untuk memperbarui status bahaya; semuanya diproses otomatis di tingkat server.

### 3.2 Pemetaan WebGIS Terintegrasi (Kawasan Rawan Bencana)
- **REQ-2.1:** Platform harus memvisualisasikan poligon Kawasan Rawan Bencana (KRB 1, 2, 3) Gunung Kelud dengan warna kontras (merah, kuning).
- **REQ-2.2:** Menampilkan alur lahar dingin Kali Ngobo (tipe *LINESTRING*).
- **REQ-2.3:** Menampilkan titik (*POINT*) lokasi posko pengungsian, titik kumpul, dan fasilitas kesehatan.

### 3.3 Inovasi Spatial AI: NLP Text Summarizer
- **REQ-3.1:** Mesin AI berbasis NLP di *backend* harus mengonversi data teknis mentah dari API (seperti amplitudo tremor, curah hujan harian dalam mm/jam) menjadi kalimat imbauan peringatan dini yang ringkas dan bahasa yang mudah dipahami warga.
- **REQ-3.2:** Imbauan ditampilkan seketika pada *Alert Banner* di antarmuka publik.

### 3.4 Inovasi Spatial AI: Semantic Search & Rute Cerdas
- **REQ-4.1:** Fitur pencarian pintar (*Semantic Search*) menggunakan basis data PostgreSQL dengan ekstensi **pgvector**.
- **REQ-4.2:** Sistem harus memproses pertanyaan bahasa sehari-hari warga (misal: "lari ke mana kalau ada lahar?") menjadi *embedding vector* dan mencocokkannya dengan posko pengungsian terdekat.
- **REQ-4.3:** Menggunakan **pgRouting** (Algoritma Dijkstra) pada basis data PostGIS untuk mengkalkulasi dan menampilkan rute evakuasi terpendek di peta yang menghindari zona bahaya (Dynamic Hazard Avoidance).

### 3.5 Antarmuka Publik
- **REQ-5.1:** Peta interaktif skala penuh (*full screen*) menggunakan SvelteKit dan Leaflet.
- **REQ-5.2:** *Interactive Pop-up* yang memuat nama fasilitas, kapasitas pengungsi, dan tombol rute evakuasi.

## 4. Kebutuhan Non-Fungsional (NFR)
- **NFR-1 (Kinerja):** Pencarian rute (*pgRouting*) dan *Semantic Search* harus merespons di bawah hitungan milidetik meskipun menanggung *traffic* warga yang tinggi saat darurat.
- **NFR-2 (Keandalan):** *Backend* Golang harus mampu menangani kegagalan koneksi eksternal (PVMBG/BMKG API) secara anggun (mitigasi *error server*).

## 5. Rencana Fase Pengerjaan (Agile/Scrum - 4 Sprint)
- **Sprint 1 (15-18 Sept):** Analisis kebutuhan dan perancangan basis data spasial PostgreSQL/PostGIS.
- **Sprint 2 (19-22 Sept):** Pembangunan *backend* Golang (Zero-Admin) dan Modul AI NLP Summarizer.
- **Sprint 3 (23-26 Sept):** Pembangunan *frontend* SvelteKit dan integrasi *Semantic Search* (pgvector).
- **Sprint 4 (27-30 Sept):** Pengujian Teknis (*Black-Box*) dan Evaluasi Pengguna (TAM).

---
*Status: SELARAS DENGAN PROPOSAL AKADEMIK*
