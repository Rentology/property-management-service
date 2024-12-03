package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	service2 "property-managment-service/internal/image/service"
	"property-managment-service/internal/models"
)

type imageRepository struct {
	Db *sqlx.DB
}

func NewImageRepository(db *sqlx.DB) service2.ImageRepository {
	return &imageRepository{Db: db}
}

func (r *imageRepository) SaveImage(ctx context.Context, image *models.Image) (*models.Image, error) {
	const op = "imageRepository.SaveImage"
	query := `INSERT INTO properties_images (property_id, image_url) VALUES ($1, $2) RETURNING *`
	if err := r.Db.QueryRowxContext(ctx, query, image.PropertyId, image.ImageUrl).StructScan(image); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return image, nil
}

func (r *imageRepository) GetImage(ctx context.Context, id int64) (*models.Image, error) {
	const op = "imageRepository.GetImage"
	query := `SELECT * FROM properties_images WHERE ID = $1`
	image := &models.Image{}
	if err := r.Db.QueryRowxContext(ctx, query, id).StructScan(image); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return image, nil
}

func (r *imageRepository) GetImagesByPropertyID(ctx context.Context, propertyID int64) ([]models.Image, error) {
	const op = "imageRepository.GetImagesByPropertyID"
	query := `SELECT * FROM properties_images WHERE property_id = $1`
	rows, err := r.Db.QueryContext(ctx, query, propertyID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var image models.Image
		err := rows.Scan(&image.Id, &image.PropertyId, &image.ImageUrl)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		images = append(images, image)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return images, nil
}

func (r *imageRepository) SaveImageWithTx(ctx context.Context, image *models.Image, tx *sqlx.Tx) (*models.Image, error) {
	const op = "imageRepository.SaveImage"
	query := `INSERT INTO properties_images (property_id, image_url) VALUES ($1, $2) RETURNING *`

	// Используем переданную транзакцию
	if err := tx.QueryRowxContext(ctx, query, image.PropertyId, image.ImageUrl).StructScan(image); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return image, nil
}

func (r *imageRepository) DeleteWithTx(ctx context.Context, id int64, tx *sqlx.Tx) error {
	const op = "imageRepository.DeleteWithTx"
	query := `DELETE FROM properties_images WHERE id = $1 RETURNING id`

	var deletedID int64
	// Выполняем DELETE и возвращаем id удалённой записи
	err := tx.QueryRowxContext(ctx, query, id).Scan(&deletedID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
