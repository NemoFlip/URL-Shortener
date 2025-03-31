package app

import (
	_ "RESTProject/docs"
	"RESTProject/internal/config"
	mwLogger "RESTProject/internal/delivery/middleware/logger"
	"RESTProject/internal/delivery/routes"
	"RESTProject/internal/lib/logger/sl"
	"RESTProject/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
)

// @title REST API Service
// @description REST API project ready for prod
// @host localhost:8080
// @BasePath /
func Run(cfg *config.Config, logger *slog.Logger) {
	storage, err := postgres.NewStorage(cfg.DataSourceName)
	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	logger.Info("Storage initialized")

	_ = storage

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(mwLogger.New(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)



	r.Get("/swagger/*", httpSwagger.WrapHandler)

	routes.InitRouting(logger, storage, r, cfg)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	logger.Info("starting server", slog.String("address", cfg.Address))

	if err = srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}
	logger.Error("server stopped")
}
