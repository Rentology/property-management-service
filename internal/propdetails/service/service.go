package service

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"property-managment-service/internal/models"
	"property-managment-service/internal/propdetails/delivery/http"
)

type PropertyDetailsRepository interface {
	Create(ctx context.Context, details *models.PropertyDetails) (*models.PropertyDetails, error)
	GetById(ctx context.Context, id int64) (*models.PropertyDetails, error)
	Update(ctx context.Context, details *models.PropertyDetails) (*models.PropertyDetails, error)
	Delete(ctx context.Context, id int64) (int64, error)
	SaveWithTx(ctx context.Context, details *models.PropertyDetails, tx *sqlx.Tx) error
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

func (s *propertyDetailsService) GetById(ctx context.Context, id int64) (*models.PropertyDetails, error) {
	details, err := s.propertyDetailsRepository.GetById(ctx, id)
	s.log.Info("GetById", "details", details)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (s *propertyDetailsService) Update(ctx context.Context, details *models.PropertyDetails) (*models.PropertyDetails, error) {
	details, err := s.propertyDetailsRepository.Update(ctx, details)
	s.log.Info("Update", "updated details", details)
	if err != nil {
		return nil, err
	}
	return details, err
}

func (s *propertyDetailsService) Delete(ctx context.Context, id int64) (int64, error) {
	_, err := s.propertyDetailsRepository.Delete(ctx, id)
	s.log.Info("Delete", "details id", id)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (s *propertyDetailsService) SaveWithTx(ctx context.Context, details *models.PropertyDetails, tx *sqlx.Tx) error {
	return s.propertyDetailsRepository.SaveWithTx(ctx, details, tx)
}
