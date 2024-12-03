package middleware

import (
	"log/slog"
	"property-managment-service/internal/config"
)

type MiddlewareManager struct {
	log *slog.Logger
	cfg *config.Config
}

func NewMiddlewareManager(log *slog.Logger, cfg *config.Config) *MiddlewareManager {
	return &MiddlewareManager{log: log, cfg: cfg}
}
