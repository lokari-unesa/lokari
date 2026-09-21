# UI/UX & Design System - LOKARI

## 1. Pendekatan Desain: "Active and Simplified UI/UX"
Berdasarkan proposal akademik, pengguna utama LOKARI adalah masyarakat Desa Jarak yang mayoritas awam teknologi (gaptek). Oleh karena itu, antarmuka tidak boleh membingungkan. Desain difokuskan pada keaktifan sistem (menampilkan data darurat otomatis) dengan tampilan sesederhana mungkin (Simplified UI/UX).

## 2. Kriteria Usability Goals (Sesuai Konsep TAM/UTAUT)
Platform dievaluasi berdasarkan enam prinsip kenyamanan:
1. **Kemudahan Dipelajari (Learnability):** Tata letak alami. Geser peta dan klik ikon tanpa perlu membaca petunjuk.
2. **Efisiensi Penggunaan (Efficiency):** Informasi posko atau rute ditemukan maksimal dalam 1-2 klik.
3. **Efektivitas Informasi (Effectiveness):** Peta menjawab kebutuhan darurat (visualisasi bahaya seketika).
4. **Kerapian Visual & Kontras Warna (Clarity/Visibility):** 
   - 🔴 **Merah:** Kawasan Rawan Bencana (KRB 3) & Jalur Lahar Aktif.
   - 🟡 **Kuning:** Zona Waspada (KRB 2).
   - 🟢 **Hijau:** Lokasi Posko Pengungsian yang aman.
5. **Aksesibilitas Multi-Perangkat (Responsiveness):** Berbasis SvelteKit agar responsif di ponsel, tempat warga biasa mengakses internet.
6. **Kejelasan Komponen Grafis:** Simbol *marker* peta sangat intuitif (Ikon tenda/balai desa untuk posko, ikon palang untuk kesehatan).

## 3. Komponen Antarmuka Publik (3 Pilar Utama)
Proposal mengamanatkan 3 komponen tampilan utama untuk halaman *front-end*:

### 3.1 Halaman Utama (Landing Page / Peta Penuh)
- Seluruh layar (*full screen*) didominasi oleh kanvas peta interaktif Leaflet.
- Bagian atas peta dilengkapi dengan Bilah Pencarian Pintar (*Semantic Search*) mengambang. Fitur ini dirancang seperti kolom *chat* agar warga bisa mengetik menggunakan bahasa sehari-hari.

### 3.2 Jendela Informasi Pop-up (Interactive Pop-up)
- Muncul ketika ikon posko/titik kumpul diklik.
- Memuat foto kondisi lokasi, daya tampung pengungsi, dan tombol besar "Arahkan Rute Evakuasi" (*pgRouting*).

### 3.3 Panel Imbauan Darurat (Alert Banner)
- Terletak di atas peta atau muncul sebagai *pop-up toast*.
- Memuat hasil ringkasan *AI NLP Summarizer* Golang. Contoh: "Status Gunung Waspada, jauhi aliran sungai."
- Warna banner akan berubah (Hijau/Kuning/Merah) bergantung dari level peringatan.

## 4. Resolusi Konflik Konsep Awal
*Catatan: Pada kesepakatan tim sebelumnya, terdapat rencana pembuatan Hero Carousel Youtube. Namun, karena proposal akademik mensyaratkan peta sebagai Landing Page utama demi efisiensi saat krisis, jika Carousel ingin dipertahankan, ia akan diposisikan sebagai halaman "Tentang/Edukasi" terpisah, atau sebuah layar *intro* (splash screen) statis sebelum memasuki peta.*

---
*Status: SELARAS DENGAN PROPOSAL AKADEMIK*
