package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"property-managment-service/internal/models"
	"property-managment-service/internal/property/service"
)

type propertyRepository struct {
	Db *sqlx.DB
}

func NewPropertyRepository(db *sqlx.DB) service.PropertyRepository {
	return &propertyRepository{Db: db}
}

func (r *propertyRepository) Create(ctx context.Context, property *models.Property) (*models.Property, error) {
	const op = "propertyRepository.create"
	query := `INSERT INTO properties (owner_id, title, location, price, property_type, rental_type, max_guests, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	if err := r.Db.QueryRowxContext(ctx, query, &property.OwnerId, &property.Title, &property.Location, &property.Price,
		&property.PropertyType, &property.RentalType, &property.MaxGuests, &property.CreatedAt).StructScan(property); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return property, nil
}
