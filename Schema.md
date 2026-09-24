# Skema Basis Data (PostgreSQL / PostGIS) - LOKARI

## 1. Ketentuan Spasial (SRID: 4326 WGS 84)
Berbeda dengan skema awal yang menggunakan NoSQL MongoDB, sesuai proposal sistem harus dibangun di atas relasional **PostgreSQL** dengan ekstensi **PostGIS**, **pgRouting**, dan **pgvector**.

## 2. Struktur Tabel Relasional

### 2.1 Tabel: `potensi_bencana` (Data Fasilitas & Posko)
Menyimpan lokasi titik kumpul, posko pengungsian, dan fasilitas kesehatan.
```sql
CREATE TABLE potensi_bencana (
    id_potensi UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_objek VARCHAR(255) NOT NULL,
    kategori VARCHAR(100) NOT NULL, -- 'Posko_Pengungsian', 'Fasilitas_Kesehatan'
    tingkat_risiko VARCHAR(50),     -- 'Aman', 'Rawan'
    deskripsi TEXT,
    alamat_dusun VARCHAR(255),
    kapasitas_orang INTEGER,
    foto_lokasi VARCHAR(255),
    kontak_darurat VARCHAR(50),
    geometri GEOMETRY(POINT, 4326), -- PostGIS Point
    embedding VECTOR(1024)          -- pgvector, vektor dari Cohere embed-multilingual-v3.0 (1024 dimensi)
);

-- Indexing untuk kecepatan pencarian spasial dan semantik
CREATE INDEX idx_potensi_geom ON potensi_bencana USING GiST (geometri);
CREATE INDEX idx_potensi_embed ON potensi_bencana USING hnsw (embedding vector_cosine_ops);
```

### 2.2 Tabel: `kawasan_rawan` (Poligon KRB Gunung Kelud)
```sql
CREATE TABLE kawasan_rawan (
    id_krb SERIAL PRIMARY KEY,
    zona VARCHAR(50) NOT NULL,      -- 'KRB 1', 'KRB 2', 'KRB 3'
    deskripsi_bahaya TEXT,
    geometri GEOMETRY(POLYGON, 4326) -- PostGIS Polygon
);
CREATE INDEX idx_krb_geom ON kawasan_rawan USING GiST (geometri);
```

### 2.3 Tabel: `jaringan_jalan` (Untuk pgRouting Evakuasi)
```sql
CREATE TABLE jaringan_jalan (
    id INTEGER PRIMARY KEY,
    source INTEGER,                 -- Titik awal (Node)
    target INTEGER,                 -- Titik akhir (Node)
    cost FLOAT8,                    -- Bobot jarak normal
    reverse_cost FLOAT8,            -- Bobot arah sebaliknya
    is_safe BOOLEAN DEFAULT TRUE,   -- Jika memotong zona KRB 3, set FALSE & cost = 999999
    geometri GEOMETRY(LINESTRING, 4326)
);
```

### 2.4 Tabel: `log_peringatan_dini` (Hasil NLP Summarizer)
Tabel ini diisi secara otomatis oleh sistem *Zero-Admin* Golang.
```sql
CREATE TABLE log_peringatan_dini (
    id SERIAL PRIMARY KEY,
    sumber_api VARCHAR(50),         -- 'PVMBG', 'BMKG'
    data_mentah JSONB,              -- Data teknis (tremor, curah hujan)
    pesan_ringkasan_nlp TEXT,       -- Hasil rangkuman AI untuk disajikan ke publik
    tingkat_bahaya VARCHAR(20),     -- 'Aman', 'Waspada', 'Awas'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 3. Kontrak REST API (Golang)

### 3.1 GET `/api/v1/peringatan-dini`
- Mengembalikan *Alert Banner* terakhir yang telah dirangkum AI.
```json
{
  "status": "WASPADA",
  "pesan": "Curah hujan tinggi di puncak Gunung Kelud berpotensi memicu lahar dingin di sepanjang Kali Ngobo. Warga Dusun Simbar Kidul diminta menjauhi bantaran sungai."
}
```

### 3.2 POST `/api/v1/semantic-search`
- Menerima bahasa sehari-hari dan menggunakan `pgvector` untuk mencari lokasi.
- **Request:** `{"query": "kalau lahar turun lari ke mana"}`
- **Response:** Mengembalikan GeoJSON titik posko terdekat.

### 3.3 GET `/api/v1/rute-evakuasi?lat=-7.1&lng=110.2`
- Menjalankan kueri *pgRouting* Dijkstra di basis data.
- **Response:** GeoJSON *LineString* rute teraman menuju posko pengungsian.

---
*Status: SELARAS DENGAN PROPOSAL AKADEMIK*
