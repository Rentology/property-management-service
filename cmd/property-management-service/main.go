package main

import (
	"github.com/jmoiron/sqlx"
	"log/slog"
	"os"
	"property-managment-service/internal/config"
	"property-managment-service/internal/server"
	"property-managment-service/pkg/db"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.LoadConfig()
	log := setupLogger(cfg.App.Env)

	log.Info("starting property-management-service", slog.String("env", cfg.App.Env))
	log.Debug("debug messages are enabled")

	psqlDB, err := db.NewPsqlDB(cfg)

	if err != nil {
		log.Error("failed to connect to postgresql", "error", err)
		os.Exit(1)
	} else {
		log.Info("Postgres connected", "status", psqlDB.Stats())
	}

	defer func(psqlDB *sqlx.DB) {
		err := psqlDB.Close()
		if err != nil {
			log.Error("failed to close connection", "error", err)
			os.Exit(1)
		}
	}(psqlDB)

	s := server.NewServer(cfg, psqlDB, log)
	if err = s.Run(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}
