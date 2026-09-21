package services

import (
	"log/slog"

	"github.com/dv1x3r/chinook-go/internal/config"
	"github.com/dv1x3r/chinook-go/internal/lib/db"
)

type Service struct {
	logger *slog.Logger
	store  db.Store
	cfg    config.Config
}

func NewService(logger *slog.Logger, store db.Store, cfg config.Config) *Service {
	return &Service{
		logger: logger,
		store:  store,
		cfg:    cfg,
	}
}
