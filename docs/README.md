# Dokumentasi Teknis LOKARI

Dokumen di folder ini menggambarkan sistem **sebagaimana diimplementasikan**.
Dokumen proposal akademik di root (`PRD.md`, `Design.md`, `General.md`,
`Rules.md`) sudah ditulis ulang mengikuti realita kode — baca keduanya
tanpa takut bertentangan.

| Dokumen | Isi |
|---|---|
| [`Architecture.md`](Architecture.md) | Diagram sistem, stack frontend/backend, worker Zero-Admin, AI, Web Push, deployment & port |
| [`Schema.md`](Schema.md) | Skema database aktual (6 tabel, 3 ekstensi), index, tabel proposal yang tidak dibuat + alasannya |
| [`API.md`](API.md) | Kontrak endpoint REST, aturan routing `/api/*`, rate limit, env yang dibaca |

Operasional harian (migrasi, seed, cron manual, VAPID, deploy): `README.md` di root.
