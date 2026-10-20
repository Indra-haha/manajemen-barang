package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"manajemen-barang/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(ctx, query, u.Username, u.PasswordHash).Scan(&u.ID)
	if err != nil {
		if isDuplicate(err) {
			return models.ErrUsernameTaken
		}
		return err
	}
	return nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (models.User, error) {
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`
	var u models.User
	err := r.db.QueryRow(ctx, query, username).Scan(&u.ID, &u.Username, &u.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, models.ErrNotFound
	}
	return u, err
}
