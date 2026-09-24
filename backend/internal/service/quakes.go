package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lokari/backend/internal/ai"
)

// Koordinat Gunung Kelud (pusat filter radius) & kunci state dedupe.
const (
	keludLat = -7.932
	keludLon = 112.308

	// stateKeyLastQuake menyimpan DateTime gempa terakhir yang sudah dilihat.
	stateKeyLastQuake = "bmkg_last_datetime"
)

var bmkgClient = &http.Client{Timeout: 10 * time.Second}

// bmkgQuake mewakili satu record gempa dari BMKG Open Data (TEWS).
type bmkgQuake struct {
	DateTime    string `json:"DateTime"`
	Coordinates string `json:"Coordinates"`
	Magnitude   string `json:"Magnitude"`
	Kedalaman   string `json:"Kedalaman"`
	Wilayah     string `json:"Wilayah"`
	Dirasakan   string `json:"Dirasakan"`
}

// FetchBMKGFeltQuakes menarik daftar gempa dirasakan terbaru dari BMKG,
// menyaring yang berpotensi dirasakan warga sekitar Gunung Kelud
// (magnitudo >= 3.5 + radius dinamis), lalu mengirim notifikasi untuk
// gempa yang benar-benar baru.
func (s *FetcherService) FetchBMKGFeltQuakes() {
	ctx := context.Background()
	if s.DB == nil {
		return
	}

	var candidates []bmkgQuake
	for _, url := range []string{
		"https://data.bmkg.go.id/DataMKG/TEWS/gempadirasakan.json",
		"https://data.bmkg.go.id/DataMKG/TEWS/autogempa.json",
	} {
		q, err := fetchBMKGQuakes(url)
		if err != nil {
			log.Printf("[Gempa] Gagal tarik %s: %v\n", url, err)
			continue
		}
		candidates = append(candidates, q...)
	}
	if len(candidates) == 0 {
		return
	}

	// Gempa yang sama bisa muncul di dua sumber — buang duplikat.
	seen := map[string]bool{}
	unique := make([]bmkgQuake, 0, len(candidates))
	for _, q := range candidates {
		key := q.DateTime + "|" + q.Coordinates
		if seen[key] {
			continue
		}
		seen[key] = true
		unique = append(unique, q)
	}

	lastDT := s.loadLastQuakeDT(ctx)
	maxDT := time.Time{}
	processed := 0
	for _, q := range unique {
		dt, ok := parseBMKGTime(q.DateTime)
		if !ok {
			continue
		}
		if dt.After(maxDT) {
			maxDT = dt
		}
		if !dt.After(lastDT) {
			continue // sudah pernah diproses sebelumnya
		}
		if !s.isRelevantQuake(q) {
			continue
		}
		s.saveQuakeNews(ctx, q)
		processed++
	}

	// Majukan state ke DateTime terbaru yang terlihat (termasuk gempa yang
	// tidak lolos filter) agar tidak dievaluasi ulang terus-menerus.
	if maxDT.After(lastDT) {
		if err := s.SetState(ctx, stateKeyLastQuake, maxDT.Format(time.RFC3339)); err != nil {
			log.Printf("[Gempa] Gagal simpan state last datetime: %v\n", err)
		}
	}
	if processed > 0 {
		log.Printf("[Gempa] Sesi selesai: %d gempa relevan baru diproses\n", processed)
	}
}

func (s *FetcherService) loadLastQuakeDT(ctx context.Context) time.Time {
	v, ok, err := s.GetState(ctx, stateKeyLastQuake)
	if err != nil || !ok {
		return time.Time{}
	}
	t, ok := parseBMKGTime(v)
	if !ok {
		return time.Time{}
	}
	return t
}

// isRelevantQuake: magnitudo >= 3.5 DAN jarak episenter ke Gunung Kelud
// masih dalam radius dinamis (makin besar M, makin lebar area yang dirasakan).
func (s *FetcherService) isRelevantQuake(q bmkgQuake) bool {
	m, err := parseMagnitude(q.Magnitude)
	if err != nil || m < 3.5 {
		return false
	}
	lat, lon, err := parseCoords(q.Coordinates)
	if err != nil {
		return false
	}
	return haversineKm(lat, lon, keludLat, keludLon) <= keludRadius(m)
}

func (s *FetcherService) saveQuakeNews(ctx context.Context, q bmkgQuake) {
	raw, _ := json.Marshal(map[string]string{
		"waktu":     q.DateTime,
		"koordinat": q.Coordinates,
		"magnitude": q.Magnitude,
		"kedalaman": q.Kedalaman,
		"wilayah":   q.Wilayah,
		"dirasakan": q.Dirasakan,
	})

	title, summary := fallbackQuakeText(q)
	if news, err := ai.GenerateNewsSummary(ctx, string(raw), "BMKG (Badan Meteorologi, Klimatologi, dan Geofisika)"); err == nil {
		if news.Title != "" {
			title = news.Title
		}
		if news.Summary != "" {
			summary = news.Summary
		}
	}

	// Kategori dipaksa "warning": gempa yang lolos filter sudah dipastikan
	// berpotensi dirasakan warga sekitar Gunung Kelud → layak notifikasi.
	query := `INSERT INTO kabar_kelud (kategori, judul, ringkasan, sumber) VALUES ($1, $2, $3, $4) ON CONFLICT (sumber, judul) DO NOTHING`
	tag, err := s.DB.Exec(ctx, query, "warning", title, summary, "BMKG")
	if err != nil {
		log.Printf("[Gempa] Gagal simpan berita: %v\n", err)
		return
	}
	if tag.RowsAffected() != 1 {
		return // judul+sumber sama sudah pernah ada — jangan notif ulang
	}
	log.Printf("[Gempa] Berita baru tersimpan: %s\n", title)
	go s.broadcastPush(title, summary, "warning")
}

// keludRadius memperkirakan jarak maksimum gempa yang masih berpotensi
// dirasakan warga sekitar Gunung Kelud: M4→50 km, M5→150 km, M6→350 km,
// interpolasi linear per segmen, cap 350 km.
func keludRadius(m float64) float64 {
	switch {
	case m >= 6:
		return 350
	case m >= 5:
		return 150 + (m-5)*200
	default:
		r := 50 + (m-4)*100
		if r < 50 {
			return 50
		}
		return r
	}
}

func fetchBMKGQuakes(url string) ([]bmkgQuake, error) {
	resp, err := bmkgClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status HTTP %d", resp.StatusCode)
	}
	var raw struct {
		Infogempa struct {
			Gempa json.RawMessage `json:"gempa"`
		} `json:"Infogempa"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	if len(raw.Infogempa.Gempa) > 0 && raw.Infogempa.Gempa[0] == '[' {
		var list []bmkgQuake
		if err := json.Unmarshal(raw.Infogempa.Gempa, &list); err != nil {
			return nil, err
		}
		return list, nil
	}
	var one bmkgQuake
	if err := json.Unmarshal(raw.Infogempa.Gempa, &one); err != nil {
		return nil, err
	}
	return []bmkgQuake{one}, nil
}

// parseBMKGTime toleran terhadap format DateTime BMKG (RFC3339 maupun
// "2006-01-02 15:04:05").
func parseBMKGTime(v string) (time.Time, bool) {
	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, v); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func parseMagnitude(v string) (float64, error) {
	return strconv.ParseFloat(strings.Replace(strings.TrimSpace(v), ",", ".", 1), 64)
}

func parseCoords(v string) (float64, float64, error) {
	parts := strings.Split(v, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("format koordinat tidak dikenal: %s", v)
	}
	lat, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	lon, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("format koordinat tidak dikenal: %s", v)
	}
	return lat, lon, nil
}

func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371.0
	toRad := func(d float64) float64 { return d * math.Pi / 180 }
	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	return 2 * r * math.Asin(math.Sqrt(a))
}

func fallbackQuakeText(q bmkgQuake) (string, string) {
	title := fmt.Sprintf("Gempa M%s di %s", q.Magnitude, q.Wilayah)
	summary := q.Dirasakan
	if summary == "" {
		summary = fmt.Sprintf("Terjadi gempa berkekuatan M%s pada %s (pusat %s, kedalaman %s). Tetap tenang dan ikuti arahan BPBD setempat.", q.Magnitude, q.DateTime, q.Wilayah, q.Kedalaman)
	}
	return title, summary
}