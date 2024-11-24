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
	GetById(ctx context.Context, id int64) (*models.Property, error)
	GetByOwnerId(ctx context.Context, id int64) ([]*models.Property, error)
	Update(ctx context.Context, property *models.Property) (*models.Property, error)
	Delete(ctx context.Context, id int64) (int64, error)
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

func (s *propertyService) GetById(ctx context.Context, id int64) (*models.Property, error) {
	property, err := s.propertyRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return property, nil
}

func (s *propertyService) GetByOwnerId(ctx context.Context, id int64) ([]*models.Property, error) {
	properties, err := s.propertyRepo.GetByOwnerId(ctx, id)
	if err != nil {
		return nil, err
	}
	return properties, nil
}

func (s *propertyService) Delete(ctx context.Context, id int64) (int64, error) {
	deleteId, err := s.propertyRepo.Delete(ctx, id)
	if err != nil {
		return 0, err
	}
	return deleteId, nil
}

func (s *propertyService) Update(ctx context.Context, property *models.Property) (*models.Property, error) {
	property, err := s.propertyRepo.Update(ctx, property)
	if err != nil {
		return nil, err
	}
	return property, nil
}
