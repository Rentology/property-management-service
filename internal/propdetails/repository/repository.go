package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"property-managment-service/internal/models"
	"property-managment-service/internal/propdetails/service"
)

type propDetailsRepository struct {
	Db *sqlx.DB
}

func NewPropDetailsRepository(db *sqlx.DB) service.PropertyDetailsRepository {
	return &propDetailsRepository{Db: db}
}

func (r *propDetailsRepository) Create(ctx context.Context, details *models.PropertyDetails) (*models.PropertyDetails, error) {
	const op = "propDetailsRepository.Create"
	query := `INSERT INTO property_details(property_id, floor, max_floor, area, rooms, house_creation_year, house_type, description)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	if err := r.Db.QueryRowxContext(ctx, query, details.PropertyID, details.Floor, details.MaxFloor, details.Area,
		details.Rooms, details.HouseCreationYear, details.HouseType, details.Description).StructScan(details); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return details, nil

}

func (r *propDetailsRepository) GetById(ctx context.Context, id int64) (*models.PropertyDetails, error) {
	const op = "propDetailsRepository.GetById"
	query := `SELECT * FROM property_details WHERE property_id = $1`
	details := &models.PropertyDetails{}
	if err := r.Db.QueryRowxContext(ctx, query, id).StructScan(details); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return details, nil
}

func (r *propDetailsRepository) Delete(ctx context.Context, id int64) (int64, error) {
	const op = "propDetailsRepository.GetById"
	var deletedId int64
	query := `DELETE FROM property_details WHERE property_id = $1 RETURNING property_id`
	if err := r.Db.QueryRowxContext(ctx, query, id).Scan(&deletedId); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return deletedId, nil
}

func (r *propDetailsRepository) Update(ctx context.Context, details *models.PropertyDetails) (*models.PropertyDetails, error) {
	const op = "propDetailsRepository.Update"
	query := `UPDATE property_details 
              SET floor = $1, max_floor = $2, area = $3, rooms = $4, 
                  house_creation_year = $5, house_type = $6, description = $7	 
              WHERE property_id = $8 RETURNING *`

	if err := r.Db.QueryRowxContext(ctx, query,
		details.Floor, details.MaxFloor, details.Area, details.Rooms,
		details.HouseCreationYear, details.HouseType, details.PropertyID).StructScan(details); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return details, nil
}

func (r *propDetailsRepository) SaveWithTx(ctx context.Context, details *models.PropertyDetails, tx *sqlx.Tx) error {
	query := `INSERT INTO property_details(property_id, floor, max_floor, area, rooms, house_creation_year, house_type, description)
			  VALUES (:property_id, :floor, :max_floor, :area, :rooms, :house_creation_year, :house_type, :description)`

	_, err := tx.NamedExecContext(ctx, query, details)
	if err != nil {
		return fmt.Errorf("failed to insert property details: %w", err)
	}
	return nil
}
