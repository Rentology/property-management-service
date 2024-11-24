package service

import (
	"context"
	"log/slog"
	"property-managment-service/internal/models"
	"property-managment-service/internal/property/delivery/http"
	"property-managment-service/pkg/utils"
	"time"
)

type PropertyRepository interface {
	Create(ctx context.Context, property *models.Property) (*models.Property, error)
}

type propertyService struct {
	log          *slog.Logger
	propertyRepo PropertyRepository
}

func NewPropertyService(propertyRepo PropertyRepository, log *slog.Logger) http.PropertyService {
	return &propertyService{log, propertyRepo}
}

func (s *propertyService) Create(ctx context.Context, property *models.Property) (*models.Property, error) {

	property.CreatedAt = time.Now().Format("2006-01-2")
	property, err := s.propertyRepo.Create(ctx, property)
	if err != nil {
		return nil, err
	}
	formattedDate, err := utils.ParseDate(&property.CreatedAt)

	if err == nil {
		property.CreatedAt = formattedDate
	}
	return property, nil
}
