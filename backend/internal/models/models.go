package models

import (
	"time"
)

type Pengguna struct {
	IDPengguna   string    `json:"id_pengguna"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type KategoriLayer struct {
	IDKategori   int    `json:"id_kategori"`
	NamaKategori string `json:"nama_kategori"`
	IkonMarker   string `json:"ikon_marker"`
}

type PotensiBencana struct {
	IDPotensi      string    `json:"id_potensi"`
	NamaObjek      string    `json:"nama_objek"`
	Kategori       string    `json:"kategori"`
	TingkatRisiko  *string   `json:"tingkat_risiko,omitempty"`
	Deskripsi      *string   `json:"deskripsi,omitempty"`
	AlamatDusun    *string   `json:"alamat_dusun,omitempty"`
	KapasitasOrang *int      `json:"kapasitas_orang,omitempty"`
	Geometri       string    `json:"geometri"` // We will fetch GeoJSON string representing the geometry
	FotoLokasi     *string   `json:"foto_lokasi,omitempty"`
	KontakDarurat  *string   `json:"kontak_darurat,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type LogUpdate struct {
	IDLog       string    `json:"id_log"`
	SumberAPI   string    `json:"sumber_api"`
	StatusTarik string    `json:"status_tarik"`
	WaktuUpdate time.Time `json:"waktu_update"`
}
