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

// stateKeyKeludStatus menyimpan status terakhir Gunung Kelud di monitor_state.
const stateKeyKeludStatus = "kelud_status"

var magmaHTTP = &http.Client{Timeout: 15 * time.Second}

var (
	magmaRowRe   = regexp.MustCompile(`(?s)<tr[^>]*>(.*?)</tr>`)
	magmaLevelRe = regexp.MustCompile(`(?i)\bLevel\s*(I|II|III|IV)\s\((Normal|Waspada|Siaga|Awas)\)`)
	magmaTagRe   = regexp.MustCompile(`<[^>]+>`)
)

// FetchKeludStatus mengambil tingkat aktivitas Gunung Kelud dari halaman
// MAGMA PVMBG (magma.esdm.go.id) lalu mengirim notifikasi bila status
// berubah — naik ataupun turun (Level I…IV).
func (s *FetcherService) FetchKeludStatus() {
	ctx := context.Background()
	if s.DB == nil {
		return
	}

	level, err := fetchKeludLevel()
	if err != nil {
		log.Printf("[Kelud] Gagal mengambil status: %v\n", err)
		return
	}

	last, ok, err := s.GetState(ctx, stateKeyKeludStatus)
	if err != nil {
		log.Printf("[Kelud] Gagal membaca state: %v\n", err)
		return
	}
	if ok && last == level {
		return // tidak ada perubahan status
	}

	s.saveKeludStatusNews(ctx, level, last)
	if err := s.SetState(ctx, stateKeyKeludStatus, level); err != nil {
		log.Printf("[Kelud] Gagal simpan state: %v\n", err)
	}
}

func (s *FetcherService) saveKeludStatusNews(ctx context.Context, level, last string) {
	title := "Status Gunung Kelud: " + level
	summary := fmt.Sprintf("PVMBG mengumumkan perubahan tingkat aktivitas Gunung Kelud. Pantau terus informasi resmi dan ikuti arahan BPBD setempat.")
	if last != "" {
		summary = fmt.Sprintf("Tingkat aktivitas Gunung Kelud berubah: %s → %s. Pantau terus informasi resmi PVMBG MAGMA dan ikuti arahan BPBD setempat.", last, level)
	}

	query := `INSERT INTO kabar_kelud (kategori, judul, ringkasan, sumber) VALUES ($1, $2, $3, $4) ON CONFLICT (sumber, judul) DO NOTHING`
	tag, err := s.DB.Exec(ctx, query, "volcano", title, summary, "MAGMA")
	if err != nil {
		log.Printf("[Kelud] Gagal simpan berita: %v\n", err)
		return
	}
	if tag.RowsAffected() != 1 {
		return // berita status ini sudah pernah dibuat
	}
	log.Printf("[Kelud] Status berubah: %s\n", level)
	go s.broadcastPush(title, summary, "volcano")
}

// fetchKeludLevel mem-parse level Gunung Kelud dari halaman MAGMA:
// baris berisi "Level II (Waspada)" menetapkan level aktif, dan baris
// gunung berikutnya yang mengandung "Kelud" = level milik Kelud.
func fetchKeludLevel() (string, error) {
	resp, err := magmaHTTP.Get("https://magma.esdm.go.id/v1/gunung-api/tingkat-aktivitas")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	current := ""
	for _, m := range magmaRowRe.FindAllStringSubmatch(string(body), -1) {
		row := strings.ReplaceAll(m[1], "\u00a0", " ")
		plain := strings.Join(strings.Fields(magmaTagRe.ReplaceAllString(row, "")), " ")
		if lm := magmaLevelRe.FindStringSubmatch(plain); lm != nil {
			current = "Level " + lm[1] + " (" + capitalizeFirst(lm[2]) + ")"
			continue
		}
		if current != "" && strings.Contains(plain, "Kelud") {
			return current, nil
		}
	}
	return "", fmt.Errorf("baris Gunung Kelud tidak ditemukan di halaman MAGMA")
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}