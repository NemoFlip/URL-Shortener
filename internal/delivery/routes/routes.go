package routes

import (
	"RESTProject/internal/config"
	del "RESTProject/internal/delivery/handlers/url/delete"
	"RESTProject/internal/delivery/handlers/url/redirect"
	"RESTProject/internal/delivery/handlers/url/save"
	"RESTProject/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func InitRouting(logger *slog.Logger, s *postgres.Storage, r chi.Router, cfg *config.Config) {
	r.Route("/url", func(router chi.Router) {
		router.Use(middleware.BasicAuth("RESTProject", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		router.Post("/", save.New(logger, s))
		router.Delete("/{alias}", del.New(logger, s))


	})
	r.Get("/{alias}", redirect.New(logger, s))
}
