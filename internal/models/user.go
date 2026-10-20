package models

import (
	"errors"
	"strings"
)

var (
	ErrUnauthorized  = errors.New("username atau password salah")
	ErrUsernameTaken = errors.New("username sudah digunakan")
	ErrInvalidAkun   = errors.New("username atau password tidak valid")
)

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	r.Username = strings.TrimSpace(r.Username)
	r.Password = strings.TrimSpace(r.Password)

	if len(r.Username) < 3 || len(r.Username) > 50 || len(r.Password) < 6 || len(r.Password) > 72 {
		return ErrInvalidAkun
	}
	return nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	r.Username = strings.TrimSpace(r.Username)
	r.Password = strings.TrimSpace(r.Password)

	if r.Username == "" || r.Password == "" {
		return ErrUnauthorized
	}
	return nil
}
