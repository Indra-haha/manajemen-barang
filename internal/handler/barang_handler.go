package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"manajemen-barang/internal/models"
)

type barangStore interface {
	Create(ctx context.Context, b *models.Barang) error
	GetAll(ctx context.Context, f models.BarangFilter) ([]models.Barang, error)
	GetByID(ctx context.Context, id int64) (models.Barang, error)
	Update(ctx context.Context, b *models.Barang) error
	Delete(ctx context.Context, id int64) error
	SesuaikanStok(ctx context.Context, id int64, delta int) error
}

type BarangHandler struct {
	repo barangStore
}

func NewBarangHandler(repo barangStore) *BarangHandler {
	return &BarangHandler{repo: repo}
}

func (h *BarangHandler) Create(c *gin.Context) {
	b, ok := parseBodyBarang(c)
	if !ok {
		return
	}
	if err := h.repo.Create(c.Request.Context(), &b); err != nil {
		respondModelError(c, err)
		return
	}
	c.JSON(http.StatusCreated, b)
}

func (h *BarangHandler) List(c *gin.Context) {
	f := models.BarangFilter{
		Kategori: strings.TrimSpace(c.Query("kategori")),
		Search:   strings.TrimSpace(c.Query("q")),
	}
	items, err := h.repo.GetAll(c.Request.Context(), f)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "gagal mengambil data")
		return
	}
	if items == nil {
		items = []models.Barang{}
	}
	c.JSON(http.StatusOK, items)
}

func (h *BarangHandler) GetById(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	b, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		respondModelError(c, err)
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *BarangHandler) UpdateById(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	b, ok := parseBodyBarang(c)
	if !ok {
		return
	}
	b.ID = id
	if err := h.repo.Update(c.Request.Context(), &b); err != nil {
		respondModelError(c, err)
		return
	}
	updated, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		respondModelError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *BarangHandler) DeleteById(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		respondModelError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// SesuaikanStok mengubah stok secara atomik. Delta positif menambah, negatif mengurangi.
func (h *BarangHandler) SesuaikanStok(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req models.StokRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "body tidak valid")
		return
	}
	if err := req.Validate(); err != nil {
		respondModelError(c, err)
		return
	}
	if err := h.repo.SesuaikanStok(c.Request.Context(), id, req.Jumlah); err != nil {
		respondModelError(c, err)
		return
	}
	b, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		respondModelError(c, err)
		return
	}
	c.JSON(http.StatusOK, b)
}

func parseBodyBarang(c *gin.Context) (models.Barang, bool) {
	var b models.Barang
	if err := c.ShouldBindJSON(&b); err != nil {
		respondError(c, http.StatusBadRequest, "body tidak valid")
		return models.Barang{}, false
	}
	if err := b.Validate(); err != nil {
		respondModelError(c, err)
		return models.Barang{}, false
	}
	return b, true
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "id tidak valid")
		return 0, false
	}
	return id, true
}
