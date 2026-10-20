package router

import (
	"github.com/gin-gonic/gin"

	"manajemen-barang/internal/handler"
)

func NewRouter(barang *handler.BarangHandler, auth *handler.AuthHandler, jwtSecret string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	api := r.Group("/api")
	{
		api.POST("/register", auth.Register)
		api.POST("/login", auth.Login)

		protected := api.Group("/barang", handler.AuthRequired(jwtSecret))
		{
			protected.GET("", barang.List)
			protected.POST("", barang.Create)
			protected.GET("/:id", barang.GetById)
			protected.PUT("/:id", barang.UpdateById)
			protected.PUT("/:id/stok", barang.SesuaikanStok)
			protected.DELETE("/:id", barang.DeleteById)
		}
	}

	return r
}
