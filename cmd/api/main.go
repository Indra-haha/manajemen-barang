package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"manajemen-barang/internal/config"
	"manajemen-barang/internal/handler"
	"manajemen-barang/internal/repository"
	"manajemen-barang/internal/router"
)

func main() {
	cfg := config.Load()

	db := repository.NewDB(cfg.Database.DSN())
	defer db.Close()

	barangRepo := repository.NewBarangRepository(db)
	userRepo := repository.NewUserRepository(db)
	barangHandler := handler.NewBarangHandler(barangRepo)
	authHandler := handler.NewAuthHandler(userRepo, cfg.JWTSecret, cfg.CookieSecure)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router.NewRouter(barangHandler, authHandler, cfg.JWTSecret),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server berjalan di port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server gagal: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("gagal graceful shutdown: %v", err)
	}
}
