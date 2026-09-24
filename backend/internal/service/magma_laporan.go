package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// FetchMagmaLaporan mengambil laporan harian aktivitas Gunung Kelud dari
// halaman pencarian laporan MAGMA PVMBG lalu menyimpannya sebagai berita.
// Informasi ini bersifat rutin (bukan peringatan) sehingga tidak memicu
// notifikasi push.
func (s *FetcherService) FetchMagmaLaporan() {
	ctx := context.Background()
	if s.DB == nil {
		return
	}

	now := time.Now()
	start := now.AddDate(0, 0, -7).Format("2006-01-02")
	end := now.Format("2006-01-02")
	url := fmt.Sprintf(
		"https://magma.esdm.go.id/v1/gunung-api/laporan/search/q?code=KLD&start=%s&end=%s",
		start, end,
	)

	laporan, err := fetchMagmaLaporan(url)
	if err != nil {
		log.Printf("[Magma Laporan] Gagal menarik laporan: %v\n", err)
		return
	}

	inserted := 0
	for _, l := range laporan {
		if s.saveMagmaLaporanNews(ctx, l) {
			inserted++
		}
	}
	log.Printf("[Magma Laporan] Selesai: %d laporan baru disimpan (total %d dari MAGMA)\n", inserted, len(laporan))
}

// magmaLaporan mewakili satu laporan harian aktivitas Gunung Kelud.
type magmaLaporan struct {
	Date   string // "Rabu, 23 September 2026 - 1 hari yang lalu"
	Level  string // "Level I (Normal)"
	Author string // "Dibuat oleh ... - Rabu, ..."
	Body   string // deskripsi pengamatan visual & cuaca
	URL    string // tautan halaman detail laporan
}

// saveMagmaLaporanNews menyimpan satu laporan ke kabar_kelud. Kembali true
// hanya bila baris benar-benar baru (anti duplikat lewat UNIQUE sumber+judul).
func (s *FetcherService) saveMagmaLaporanNews(ctx context.Context, l magmaLaporan) bool {
	datePart := l.Date
	if i := strings.Index(l.Date, " - "); i > 0 {
		datePart = l.Date[:i]
	}

	judul := "Laporan Harian Gunung Kelud: " + datePart
	ringkasan := l.Level
	if l.Body != "" {
		ringkasan = l.Level + " — " + l.Body
	}

	query := `INSERT INTO kabar_kelud (kategori, judul, ringkasan, sumber) VALUES ($1, $2, $3, $4) ON CONFLICT (sumber, judul) DO NOTHING`
	tag, err := s.DB.Exec(ctx, query, "volcano", judul, ringkasan, "MAGMA")
	if err != nil {
		log.Printf("[Magma Laporan] Gagal simpan berita: %v\n", err)
		return false
	}
	if tag.RowsAffected() != 1 {
		return false // laporan tanggal ini sudah pernah disimpan
	}
	log.Printf("[Magma Laporan] Berita baru tersimpan: %s\n", judul)
	return true
}

var (
	magmaDateRe   = regexp.MustCompile(`(?s)<p class="timeline-date">(.*?)</p>`)
	magmaReportRe = regexp.MustCompile(`(?s)<p class="timeline-title">.*?<a href="([^"]*laporan/\d+[^"]*)"[^>]*>Lihat Detail</a>`)
	magmaBadgeRe  = regexp.MustCompile(`(?s)<span class="badge[^"]*">(.*?)</span>`)
	magmaAuthorRe = regexp.MustCompile(`(?s)<p class="timeline-author">(.*?)</p>`)
	magmaDescRe   = regexp.MustCompile(`(?s)<p>([^<]*)</p>`)
)

// fetchMagmaLaporan mem-parse halaman pencarian laporan MAGMA. Urutan blok
// laporan dan banner tanggal di halaman selalu selaras (banner mendahului
// laporannya), sehingga indeks keduanya bisa dipasangkan satu-satu.
func fetchMagmaLaporan(url string) ([]magmaLaporan, error) {
	resp, err := magmaHTTP.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dates := magmaDateRe.FindAllStringSubmatch(string(body), -1)
	blocks := magmaReportRe.FindAllStringSubmatch(string(body), -1)

	var out []magmaLaporan
	for i, b := range blocks {
		l := magmaLaporan{URL: b[1]}
		if i < len(dates) {
			l.Date = magmaClean(dates[i][1])
		}
		if m := magmaBadgeRe.FindStringSubmatch(b[0]); m != nil {
			l.Level = magmaClean(m[1])
		}
		if m := magmaAuthorRe.FindStringSubmatch(b[0]); m != nil {
			l.Author = magmaClean(m[1])
		}
		if m := magmaDescRe.FindStringSubmatch(b[0]); m != nil {
			l.Body = strings.TrimSpace(m[1])
		}
		out = append(out, l)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("tidak ada laporan ditemukan di halaman MAGMA")
	}
	return out, nil
}

// magmaClean menghapus tag HTML dan merapikan spasi (termasuk &nbsp;).
func magmaClean(s string) string {
	s = strings.ReplaceAll(s, "\u00a0", " ")
	return strings.Join(strings.Fields(magmaTagRe.ReplaceAllString(s, "")), " ")
}