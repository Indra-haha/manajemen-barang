package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"manajemen-barang/internal/models"
)

func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func respondModelError(c *gin.Context, err error) {
	respondError(c, httpStatusFor(err), errMessage(err))
}

func httpStatusFor(err error) int {
	switch {
	case errors.Is(err, models.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, models.ErrConflict),
		errors.Is(err, models.ErrUsernameTaken),
		errors.Is(err, models.ErrStokTidakCukup):
		return http.StatusConflict
	case errors.Is(err, models.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, models.ErrInvalid),
		errors.Is(err, models.ErrStokNegatif),
		errors.Is(err, models.ErrHargaNegatif),
		errors.Is(err, models.ErrDataTerlaluPanjang),
		errors.Is(err, models.ErrInvalidAkun):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func errMessage(err error) string {
	switch httpStatusFor(err) {
	case http.StatusInternalServerError:
		return "terjadi kesalahan pada server"
	default:
		return err.Error()
	}
}
