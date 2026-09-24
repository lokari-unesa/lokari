/**
 * Filter defensif client-side untuk pencarian AI LOKARI.
 *
 * Validasi sebenarnya ada di backend (rate limit + validasi query server-side,
 * lihat backend/internal/api/handlers/ai_handler.go). Modul ini hanya
 * mencegah request yang jelas-jelas di luar konteks evakuasi bencana.
 * Regulasi ini bisa dilewati dari client — jangan jadikan satu-satunya gerbang.
 */

// Kata yang memicu "Anomali". Dihindari kata netral sehari-hari
// (tambah/kurang/dibagi/dikali, usia/umur/jomblo) agar query sah
// seperti "posko kurang jauh dari sini" tidak kena false-positive.
const BLOCKED_WORDS = [
  "hamil", "janda", "seks", "porno", "judi", "slot", "togel", "pinjol",
  "pacar", "nikah", "jual", "beli", "harga", "promo", "diskon",
  "bokep", "mesum", "anjing", "babi", "bangsat", "tolol", "goblok", "siapa",
];

const blocklist = new RegExp(`\\b(${BLOCKED_WORDS.join("|")})\\b|(\\d+\\s*[+\\-*\\/]\\s*\\d+)`, "i");

const allowlist = /\b(posko|pengungsian|aman|selamat|masjid|mushola|musholla|msjd|mshl|sekolah|sd|smp|sma|tk|mi|mts|puskesmas|puskes|rumah sakit|rs|balai|bale|lapangan|lpngn|tempat|jalan|rute|jalur|evakuasi|lahar|gunung|kelud|bencana|darurat|terdekat|dekat|desa|dusun|lokasi|titik|kumpul|panti|warga|bantuan|jarak|plosoklaten|kediri|ngobo|simbar|kidul|gedung|kantor|apotek|klinik|bidan|polindes|polsek|koramil|kecamatan)\b/i;

export const ANOMALY_MESSAGE =
  "Sistem Mendeteksi Anomali: Pencarian AI cerdas LOKARI hanya difokuskan untuk lokasi evakuasi, fasilitas darurat, dan mitigasi bencana Gunung Kelud.";

export function isQueryAllowed(query: string): boolean {
  return !blocklist.test(query) && allowlist.test(query);
}