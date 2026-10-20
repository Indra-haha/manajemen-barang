package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"manajemen-barang/internal/models"
)

const tokenTTL = 24 * time.Hour

type userStore interface {
	Create(ctx context.Context, u *models.User) error
	GetByUsername(ctx context.Context, username string) (models.User, error)
}

type AuthHandler struct {
	repo         userStore
	jwtSecret    string
	cookieSecure bool
}

func NewAuthHandler(repo userStore, jwtSecret string, cookieSecure bool) *AuthHandler {
	if jwtSecret == "" {
		panic("JWT_SECRET kosong")
	}
	return &AuthHandler{repo: repo, jwtSecret: jwtSecret, cookieSecure: cookieSecure}
}
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "body tidak valid")
		return
	}
	if err := req.Validate(); err != nil {
		respondModelError(c, err)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	user := models.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
	}

	if err := h.repo.Create(c.Request.Context(), &user); err != nil {
		respondModelError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "body tidak valid")
		return
	}
	if err := req.Validate(); err != nil {
		respondModelError(c, err)
		return
	}

	user, err := h.repo.GetByUsername(c.Request.Context(), req.Username)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			respondModelError(c, models.ErrUnauthorized)
			return
		}
		respondModelError(c, err)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		respondModelError(c, models.ErrUnauthorized)
		return
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"iat":      now.Unix(),
		"exp":      now.Add(tokenTTL).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.jwtSecret))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	// Set token sebagai HttpOnly Cookie agar tidak bisa dibaca JavaScript (mitigasi XSS).
	// Token juga tetap dikembalikan di body untuk client non-browser (mobile/third-party).
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"token",
		token,
		int(tokenTTL.Seconds()),
		"/",
		"",
		h.cookieSecure,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
	})
}
