# Skema Basis Data (PostgreSQL / PostGIS / pgvector) — LOKARI

> Dokumen ini menggambarkan skema yang **benar-benar diterapkan** oleh
> `backend/scripts/migrate.go`. Idempoten — aman dijalankan ulang
> (`CREATE TABLE/INDEX IF NOT EXISTS`).
>
> Terakhir diselaraskan dengan kode: Oktober 2026.

## 1. Ekstensi yang dipakai

```sql
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS vector;
```

**pgRouting TIDAK dipakai.** Rute evakuasi dihitung oleh layanan eksternal
(OpenRouteService → fallback OSRM) di endpoint SvelteKit
`frontend/src/routes/api/route/+server.ts`, bukan Dijkstra di database.
Lihat `docs/Architecture.md`.

## 2. Tabel

### 2.1 `kategori_layer` — LEGACY (tidak dipakai)

```sql
CREATE TABLE IF NOT EXISTS kategori_layer (
    id_kategori SERIAL PRIMARY KEY,
    nama_kategori VARCHAR(100) NOT NULL,
    ikon_marker VARCHAR(255)
);
```

Dibuat oleh migrasi, tetapi **tidak pernah di-query** oleh backend maupun
frontend. Kategori lokasi (`Posko_Pengungsian`, `Fasilitas_Kesehatan`, …)
disimpan langsung sebagai string di `potensi_bencana.kategori`.
Biarkan apa adanya; jangan bangun fitur baru di atas tabel ini.

### 2.2 `potensi_bencana` — titik lokasi & fasilitas

Sumber data untuk peta (`GET /api/potensi`) dan pencarian semantik
(`POST /api/search`).

```sql
CREATE TABLE IF NOT EXISTS potensi_bencana (
    id_potensi UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama_objek VARCHAR(255) NOT NULL,
    kategori VARCHAR(100) NOT NULL,
    tingkat_risiko VARCHAR(50),
    deskripsi TEXT,
    alamat_dusun VARCHAR(255),
    kapasitas_orang INTEGER,
    geometri GEOMETRY,              -- umumnya POINT (SRID 4326); kolom generik, bukan GEOMETRY(POINT, 4326)
    embedding vector(1024),         -- Cohere embed-multilingual-v3.0, 1024 dimensi
    foto_lokasi VARCHAR(255),
    kontak_darurat VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

Catatan:

- `deskripsi`, `kapasitas_orang`, dan kolom opsional lain **boleh NULL**.
  Handler `GET /api/potensi` melewati baris yang `geometri`-nya NULL
  (tidak menggagalkan seluruh respons).
- `embedding` NULL = baris tidak akan muncul di hasil semantic search.
  Lengkapi lewat `go run scripts/seed.go backfill`.

### 2.3 `log_update` — riwayat penarikan data

```sql
CREATE TABLE IF NOT EXISTS log_update (
    id_log UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sumber_api VARCHAR(100) NOT NULL,   -- 'BMKG', 'NASA_EONET', …
    status_tarik VARCHAR(50) NOT NULL,  -- 'Sukses' / 'Gagal'
    waktu_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

Ditulis oleh worker Zero-Admin setiap selesai satu siklus fetch. Hanya untuk
observabilitas; tidak dibaca frontend.

### 2.4 `kabar_kelud` — berita bencana

Sumber data `GET /api/news` (20 berita terbaru) dan pemicu notifikasi push.

```sql
CREATE TABLE IF NOT EXISTS kabar_kelud (
    id_kabar UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kategori VARCHAR(50) NOT NULL,  -- 'volcano' | 'lahar' | 'evac' | 'weather' | 'warning'
    judul VARCHAR(255) NOT NULL,
    ringkasan TEXT NOT NULL,
    sumber VARCHAR(100) NOT NULL,   -- 'BMKG' | 'MAGMA' | 'NASA EONET'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_kabar_kelud_sumber_judul
    ON kabar_kelud (sumber, judul);
```

Pasangan `(sumber, judul)` unik — worker memakai
`ON CONFLICT (sumber, judul) DO NOTHING` sehingga berita yang sama tidak
pernah disimpan dua kali (anti-duplikat + anti-notifikasi ganda).

### 2.5 `push_subscriptions` — pelanggan Web Push

```sql
CREATE TABLE IF NOT EXISTS push_subscriptions (
    id SERIAL PRIMARY KEY,
    endpoint TEXT UNIQUE NOT NULL,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

Diisi via `POST /api/subscribe` dari browser warga. Entri yang mati
(HTTP 410/404 dari push service) dihapus otomatis saat broadcast berikutnya.

### 2.6 `monitor_state` — state pemantauan (sumber kebenaran status Kelud)

```sql
CREATE TABLE IF NOT EXISTS monitor_state (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

Kunci yang dipakai saat ini:

| `key`               | Isi contoh               | Ditulis oleh                              | Dibaca oleh              |
|---------------------|--------------------------|-------------------------------------------|--------------------------|
| `kelud_status`      | `Level II (Waspada)`     | Hot loop: `FetchKeludStatus`              | `GET /api/kelud/status`  |
| `bmkg_last_datetime`| `2026-09-24T10:15:00Z`  | Hot loop: `FetchBMKGFeltQuakes` (dedupe)  | Internal worker          |

## 3. Index

```sql
CREATE INDEX IF NOT EXISTS idx_potensi_geom
    ON potensi_bencana USING GiST (geometri);
CREATE INDEX IF NOT EXISTS idx_potensi_embed
    ON potensi_bencana USING hnsw (embedding vector_cosine_ops);
```

GiST mempercepat filter spasial; HNSW + cosine distance mempercepat
`ORDER BY embedding <=> $1` di semantic search.

## 4. Tabel yang TIDAK ada (sengaja tidak dibuat)

Dokumen proposal lama menyebut tiga tabel ini — semuanya **tidak
diimplementasikan** dan tidak boleh dicari di database:

| Tabel proposal | Status | Pengganti aktual |
|---|---|---|
| `kawasan_rawan` (poligon KRB 1/2/3) | ❌ tidak ada | Zona bahaya = lingkaran Leaflet hardcoded (radius 5/10/15 km) di `KeludMapView.svelte` |
| `jaringan_jalan` (topologi pgRouting) | ❌ tidak ada | Routing via ORS (`avoid_polygons`) + fallback OSRM |
| `log_peringatan_dini` (hasil NLP) | ❌ tidak ada | Ringkasan AI disimpan sebagai berita di `kabar_kelud`; riwayat fetch di `log_update` |

Alasan: tidak pernah ada digitasi topologi jaringan jalan Desa Jarak yang
layak untuk pgRouting, dan poligon KRB resmi belum tersedia sebagai data
geospasial. Lihat `docs/Architecture.md` § Routing & Zona Bahaya.

## 5. Operasional

```bash
docker compose exec backend go run scripts/migrate.go         # terapkan skema (idempoten)
docker compose exec backend go run scripts/migrate.go reset   # DESTRUKTIF: drop semua tabel + bangun ulang (konfirmasi y/N)
```

Di image produksi (tanpa toolchain Go): ganti `go run scripts/migrate.go`
dengan binary `lokari-migrate`.

Kontrak endpoint yang membaca tabel-tabel di atas: `docs/API.md`.
