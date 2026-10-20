package models

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrNotFound           = errors.New("barang tidak ditemukan")
	ErrConflict           = errors.New("kode barang sudah digunakan")
	ErrInvalid            = errors.New("data barang tidak valid")
	ErrStokNegatif        = errors.New("stok tidak boleh negatif")
	ErrHargaNegatif       = errors.New("harga tidak boleh negatif")
	ErrDataTerlaluPanjang = errors.New("data melebihi panjang maksimal")
	ErrStokTidakCukup     = errors.New("stok tidak mencukupi")
)

type Barang struct {
	ID        int64     `json:"id"`
	Kode      string    `json:"kode"`
	Nama      string    `json:"nama"`
	Kategori  string    `json:"kategori"`
	Stok      int       `json:"stok"`
	Harga     int64     `json:"harga"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate menormalkan sekaligus memvalidasi seluruh field Barang.
// Setelah Validate() sukses, model dijamin siap dikirim ke repository.
func (b *Barang) Validate() error {
	b.Kode = strings.TrimSpace(b.Kode)
	b.Nama = strings.TrimSpace(b.Nama)
	b.Kategori = strings.TrimSpace(b.Kategori)

	if b.Kode == "" || b.Nama == "" {
		return ErrInvalid
	}
	if b.Stok < 0 {
		return ErrStokNegatif
	}
	if b.Harga < 0 {
		return ErrHargaNegatif
	}
	if len(b.Kode) > 20 || len(b.Nama) > 100 || len(b.Kategori) > 50 {
		return ErrDataTerlaluPanjang
	}
	if b.Kategori == "" {
		b.Kategori = "umum"
	}
	return nil
}

type BarangFilter struct {
	Kategori string
	Search   string
}

// StokRequest adalah payload penyesuaian stok pada PUT /api/barang/:id/stok.
// Jumlah boleh positif (tambah) atau negatif (kurangi).
type StokRequest struct {
	Jumlah int `json:"jumlah"`
}

func (r *StokRequest) Validate() error {
	if r.Jumlah == 0 {
		return ErrInvalid
	}
	return nil
}
