package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"manajemen-barang/internal/models"
)

type BarangRepository struct {
	db *pgxpool.Pool
}

func NewBarangRepository(db *pgxpool.Pool) *BarangRepository {
	return &BarangRepository{db: db}
}

func (r *BarangRepository) Create(ctx context.Context, b *models.Barang) error {
	query := `INSERT INTO barang (kode, nama, kategori, stok, harga) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, b.Kode, b.Nama, b.Kategori, b.Stok, b.Harga).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if isDuplicate(err) {
			return models.ErrConflict
		}
		return err
	}
	return nil
}

func (r *BarangRepository) GetAll(ctx context.Context, f models.BarangFilter) ([]models.Barang, error) {
	query := `SELECT id, kode, nama, kategori, stok, harga, created_at, updated_at FROM barang WHERE 1=1`
	args := []any{}

	if f.Kategori != "" {
		args = append(args, f.Kategori)
		query += fmt.Sprintf(` AND kategori = $%d`, len(args))
	}
	if f.Search != "" {
		like := "%" + f.Search + "%"
		args = append(args, like, like)
		query += fmt.Sprintf(` AND (nama LIKE $%d OR kode LIKE $%d)`, len(args)-1, len(args))
	}
	query += ` ORDER BY id DESC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hasil []models.Barang
	for rows.Next() {
		var b models.Barang
		if err := rows.Scan(&b.ID, &b.Kode, &b.Nama, &b.Kategori, &b.Stok, &b.Harga, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		hasil = append(hasil, b)
	}
	return hasil, rows.Err()
}

func (r *BarangRepository) GetByID(ctx context.Context, id int64) (models.Barang, error) {
	query := `SELECT id, kode, nama, kategori, stok, harga, created_at, updated_at FROM barang WHERE id = $1`
	var b models.Barang
	err := r.db.QueryRow(ctx, query, id).
		Scan(&b.ID, &b.Kode, &b.Nama, &b.Kategori, &b.Stok, &b.Harga, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Barang{}, models.ErrNotFound
	}
	return b, err
}

func (r *BarangRepository) Update(ctx context.Context, b *models.Barang) error {
	query := `UPDATE barang SET nama = $2, kategori = $3, stok = $4, harga = $5, updated_at = now() WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, b.ID, b.Nama, b.Kategori, b.Stok, b.Harga)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

// SesuaikanStok mengubah stok secara atomik di database.
// delta positif menambah stok, delta negatif mengurangi stok.
// Pengecekan "stok tidak boleh negatif" menjadi bagian dari perintah UPDATE,
// sehingga aman dari race condition (lost update) saat banyak request bersamaan.
func (r *BarangRepository) SesuaikanStok(ctx context.Context, id int64, delta int) error {
	query := `UPDATE barang SET stok = stok + $2, updated_at = now() WHERE id = $1 AND stok + $2 >= 0`
	tag, err := r.db.Exec(ctx, query, id, delta)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		exists, err := r.exists(ctx, id)
		if err != nil {
			return err
		}
		if !exists {
			return models.ErrNotFound
		}
		return models.ErrStokTidakCukup
	}
	return nil
}

func (r *BarangRepository) exists(ctx context.Context, id int64) (bool, error) {
	var one int
	err := r.db.QueryRow(ctx, `SELECT 1 FROM barang WHERE id = $1`, id).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *BarangRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM barang WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func isDuplicate(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
