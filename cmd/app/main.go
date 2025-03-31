package main

import (
	"RESTProject/internal/app"
	"RESTProject/internal/config"
	"RESTProject/internal/lib/logger/sl"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()

	logger := sl.SetUpLogger(cfg.Env)
	logger.Info("Logger initialized", slog.String("env", cfg.Env))

	app.Run(cfg, logger)
}
