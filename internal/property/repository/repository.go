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

func (r *propertyRepository) GetById(ctx context.Context, id int64) (*models.Property, error) {
	const op = "propertyRepository.getById"
	query := `SELECT * FROM properties WHERE id = $1`
	property := &models.Property{}
	if err := r.Db.QueryRowxContext(ctx, query, id).StructScan(property); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return property, nil
}

func (r *propertyRepository) GetByOwnerId(ctx context.Context, id int64) ([]*models.Property, error) {
	const op = "propertyRepository.getByOwnerId"
	query := `SELECT * FROM properties WHERE owner_id = $1`
	rows, err := r.Db.QueryxContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	properties := []*models.Property{}
	for rows.Next() {
		property := &models.Property{}
		if err := rows.StructScan(property); err != nil {
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}
		properties = append(properties, property)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return properties, nil
}

func (r *propertyRepository) Update(ctx context.Context, property *models.Property) (*models.Property, error) {
	const op = "propertyRepository.update"
	query := `UPDATE properties 
              SET title = $1, location = $2, price = $3, property_type = $4, 
                  rental_type = $5, max_guests = $6 
              WHERE id = $7 RETURNING *`

	if err := r.Db.QueryRowxContext(ctx, query,
		property.Title, property.Location, property.Price, property.PropertyType,
		property.RentalType, property.MaxGuests, property.ID).StructScan(property); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return property, nil
}

func (r *propertyRepository) Delete(ctx context.Context, id int64) (int64, error) {
	const op = "propertyRepository.delete"
	query := `DELETE FROM properties WHERE id = $1 RETURNING id`

	var deletedID int64
	// Выполняем DELETE и возвращаем id удалённой записи
	err := r.Db.QueryRowxContext(ctx, query, id).Scan(&deletedID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return deletedID, nil
}
