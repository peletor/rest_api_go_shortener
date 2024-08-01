package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"rest_api_shortener/internal/config"
	"rest_api_shortener/internal/http-server/handlers/redirect"
	"rest_api_shortener/internal/http-server/handlers/url/save"
	"rest_api_shortener/internal/http-server/middleware/mwlogger"
	"rest_api_shortener/internal/logger/slogger"
	"rest_api_shortener/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()

	log := slogger.SetupLogger(cfg.Env)

	log.Info("Starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("Debug messages enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Error opening storage", slog.Any("Err", err))
		os.Exit(1)
	}

	log.Info("Successful opening storage")

	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("Starting server", slog.String("Address", cfg.Address))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("Failed to start server", slog.Any("Err", err))
	}

	log.Error("Server stopped", slog.String("Address", cfg.Address))
}
