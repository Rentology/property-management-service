package service

import (
	"context"
	"fmt"
	http2 "property-managment-service/internal/image/delivery/http"
	"property-managment-service/internal/models/request"
	http3 "property-managment-service/internal/propdetails/delivery/http"
	"property-managment-service/internal/property/delivery/http"
	"property-managment-service/pkg/db"
)

type propertyFormService struct {
	transactionManager     db.TransactionManager
	propertyService        http.PropertyService
	imageService           http2.ImageService
	propertyDetailsService http3.PropertyDetailsService
}

func NewPropertyFormService(
	transactionManager db.TransactionManager,
	propertyService http.PropertyService,
	imageService http2.ImageService,
	propertyDetailsService http3.PropertyDetailsService,
) http.PropertyFormService {
	return &propertyFormService{
		transactionManager:     transactionManager,
		propertyService:        propertyService,
		imageService:           imageService,
		propertyDetailsService: propertyDetailsService,
	}
}

func (s *propertyFormService) SavePropertyForm(ctx context.Context, form *request.AddPropertyRequest) error {
	tx, err := s.transactionManager.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, "SET CONSTRAINTS ALL DEFERRED")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to defer constraints: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err = s.propertyService.SaveWithTx(ctx, form.Property, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save property info: %w", err)
	}

	// Убедитесь, что ID установлен
	if form.Property.ID == 0 {
		tx.Rollback()
		return fmt.Errorf("property ID is not set after saving property")
	}

	// Установите PropertyID для деталей
	form.PropertyDetails.PropertyID = form.Property.ID

	if err := s.propertyDetailsService.SaveWithTx(ctx, form.PropertyDetails, tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save property details: %w", err)
	}

	if len(form.Images) > 0 {
		err = s.imageService.UploadImagesFromBase64(ctx, form.Images, form.Property.ID, tx)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to upload images: %w", err)
		}
	}

	return tx.Commit()
}

func (s *propertyFormService) DeletePropertyForm(ctx context.Context, propertyID int64) error {
	// Начало транзакции
	tx, err := s.transactionManager.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Удаление изображений, связанных с объектом
	err = s.imageService.DeleteImagesByPropertyId(ctx, propertyID, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete images: %w", err)
	}

	// Удаление деталей собственности
	err = s.propertyDetailsService.DeleteWithTx(ctx, propertyID, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete property details: %w", err)
	}

	// Удаление самой собственности
	err = s.propertyService.DeleteWithTx(ctx, propertyID, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete property: %w", err)
	}

	// Коммит транзакции
	return tx.Commit()
}
