package service

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// GetState membaca nilai state monitor. Bool kedua = apakah kunci sudah ada
// (false = belum pernah disimpan / database tidak tersedia).
func (s *FetcherService) GetState(ctx context.Context, key string) (string, bool, error) {
	if s.DB == nil {
		return "", false, nil
	}
	var v string
	err := s.DB.QueryRow(ctx, "SELECT value FROM monitor_state WHERE key = $1", key).Scan(&v)
	if err == pgx.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return v, true, nil
}

// SetState menulis nilai state monitor (upsert berdasarkan key).
func (s *FetcherService) SetState(ctx context.Context, key, value string) error {
	if s.DB == nil {
		return nil
	}
	_, err := s.DB.Exec(ctx, `
		INSERT INTO monitor_state (key, value) VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP
	`, key, value)
	return err
}