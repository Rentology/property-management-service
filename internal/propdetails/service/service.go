package service

import (
	"context"
	"log/slog"
	"property-managment-service/internal/models"
	"property-managment-service/internal/propdetails/delivery/http"
)

type PropertyDetailsRepository interface {
	Create(ctx context.Context, details *models.PropertyDetails) (*models.PropertyDetails, error)
	GetById(ctx context.Context, id int64) (*models.PropertyDetails, error)
}

type propertyDetailsService struct {
	propertyDetailsRepository PropertyDetailsRepository
	log                       *slog.Logger
}

func NewPropertyDetailsService(repository PropertyDetailsRepository, log *slog.Logger) http.PropertyDetailsService {
	return &propertyDetailsService{propertyDetailsRepository: repository, log: log}
}

func (s *propertyDetailsService) Create(ctx context.Context, details *models.PropertyDetails) (*models.PropertyDetails, error) {
	details, err := s.propertyDetailsRepository.Create(ctx, details)
	if err != nil {
		return nil, err
	}
	return details, nil
}
