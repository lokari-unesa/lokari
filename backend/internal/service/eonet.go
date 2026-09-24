package service

import "fmt"

// eonetRadiusKm membatasi event NASA EONET hingga radius yang masih masuk
// akal untuk warga sekitar Gunung Kelud (≈ praktis hanya mencakup Semeru).
const eonetRadiusKm = 150.0

// eonetNearKeludData menyaring event EONET menjadi yang memiliki titik
// geometry dalam radius 150 km dari Gunung Kelud, lalu membatasinya ke
// 5 event agar ringkasan AI tidak melebihi konteks.
func eonetNearKeludData(data map[string]interface{}) (interface{}, error) {
	events, ok := data["events"].([]interface{})
	if !ok || len(events) == 0 {
		return nil, fmt.Errorf("tidak ada event vulkanik aktif dari EONET")
	}
	near, err := filterNearKelud(events)
	if err != nil {
		return nil, err
	}
	if len(near) > 5 {
		near = near[:5]
	}
	return map[string]interface{}{"events": near}, nil
}

// filterNearKelud mengembalikan event yang memiliki setidaknya satu titik
// geometry dalam radius eonetRadiusKm dari Gunung Kelud. Koordinat yang
// dikirim EONET v3 adalah [longitude, latitude].
func filterNearKelud(events []interface{}) ([]interface{}, error) {
	out := make([]interface{}, 0, len(events))
	for _, ev := range events {
		m, ok := ev.(map[string]interface{})
		if !ok {
			continue
		}
		geom, ok := m["geometry"].([]interface{})
		if !ok {
			continue
		}
		near := false
		for _, g := range geom {
			point, ok := g.(map[string]interface{})
			if !ok {
				continue
			}
			coords, ok := point["coordinates"].([]interface{})
			if !ok || len(coords) < 2 {
				continue
			}
			lon, lonOK := toFloat(coords[0])
			lat, latOK := toFloat(coords[1])
			if lonOK && latOK && haversineKm(lat, lon, keludLat, keludLon) <= eonetRadiusKm {
				near = true
				break
			}
		}
		if near {
			out = append(out, ev)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("tidak ada event EONET dalam radius %g km dari Gunung Kelud", eonetRadiusKm)
	}
	return out, nil
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	}
	return 0, false
}
