package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(dsn string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(fmt.Errorf("buka koneksi postgres: %w", err))
	}
	if err := pool.Ping(context.Background()); err != nil {
		panic(fmt.Errorf("koneksi postgres gagal: %w", err))
	}
	return pool
}
