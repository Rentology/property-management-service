package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	http2 "property-managment-service/internal/image/delivery/http"
	"property-managment-service/internal/models"
	"strings"
)

type ImageRepository interface {
	SaveImage(ctx context.Context, image *models.Image) (*models.Image, error)
	SaveImageWithTx(ctx context.Context, image *models.Image, tx *sqlx.Tx) (*models.Image, error)
	GetImage(ctx context.Context, id int64) (*models.Image, error)
	GetImagesByPropertyID(ctx context.Context, propertyID int64) ([]models.Image, error)
	DeleteWithTx(ctx context.Context, id int64, tx *sqlx.Tx) error
}

type imageService struct {
	log       *slog.Logger
	imageRepo ImageRepository
}

func NewImageService(imageRepo ImageRepository, log *slog.Logger) http2.ImageService {
	return &imageService{imageRepo: imageRepo, log: log}
}

func (s *imageService) UploadImage(ctx context.Context, file *multipart.FileHeader, propertyId int64) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	homeDir := "/Users/roflandown/Desktop"

	uploadDir := filepath.Join(homeDir, "images")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm) // Создаем директорию, если её нет
		if err != nil {
			return err
		}
	}

	uniqueID := uuid.New().String()
	ext := filepath.Ext(file.Filename)
	baseName := file.Filename[:len(file.Filename)-len(ext)]
	newFileName := baseName + "_" + uniqueID + ext

	dstPath := filepath.Join(uploadDir, newFileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	if err != nil {
		return err
	}

	image := &models.Image{PropertyId: propertyId, ImageUrl: dstPath}
	_, err = s.imageRepo.SaveImage(ctx, image)
	if err != nil {
		return err
	}

	return nil
}

func (s *imageService) UploadImageFromBase64(ctx context.Context, base64Image string, propertyId int64, tx *sqlx.Tx) error {
	// Разделяем Base64 строку на часть с MIME-типом и данные
	parts := strings.SplitN(base64Image, ",", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid base64 image format")
	}

	// MIME-тип (например, data:image/png;base64)
	mimeType := parts[0]
	imageData := parts[1]

	// Декодируем Base64-строку
	decodedImage, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return fmt.Errorf("failed to decode base64 image: %w", err)
	}

	// Определяем расширение файла на основе MIME-типа
	var ext string
	switch {
	case strings.Contains(mimeType, "image/jpeg"):
		ext = ".jpg"
	case strings.Contains(mimeType, "image/png"):
		ext = ".png"
	case strings.Contains(mimeType, "image/gif"):
		ext = ".gif"
	default:
		return fmt.Errorf("unsupported image type: %s", mimeType)
	}

	// Уникальное имя файла
	uniqueID := uuid.New().String()
	newFileName := fmt.Sprintf("image_%s%s", uniqueID, ext)

	// Путь к директории для загрузки
	homeDir := "/Users/roflandown/Desktop"
	uploadDir := filepath.Join(homeDir, "images")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create upload directory: %w", err)
		}
	}

	// Путь к файлу
	dstPath := filepath.Join(uploadDir, newFileName)

	// Записываем данные в файл
	err = os.WriteFile(dstPath, decodedImage, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	// Сохраняем информацию об изображении в базе данных
	image := &models.Image{
		PropertyId: propertyId,
		ImageUrl:   dstPath,
	}
	_, err = s.imageRepo.SaveImageWithTx(ctx, image, tx)
	if err != nil {
		return fmt.Errorf("failed to save image record: %w", err)
	}

	return nil
}

func (s *imageService) UploadImagesFromBase64(ctx context.Context, base64Images []string, propertyId int64, tx *sqlx.Tx) error {
	for _, base64Image := range base64Images {
		err := s.UploadImageFromBase64(ctx, base64Image, propertyId, tx)
		if err != nil {
			return fmt.Errorf("failed to upload image: %w", err)
		}
	}
	return nil
}

func (s *imageService) GetImage(ctx context.Context, id int64) (string, *os.File, error) {
	image, err := s.imageRepo.GetImage(ctx, id)
	if err != nil {
		return "", nil, err
	}

	// Открываем файл
	file, err := os.Open(image.ImageUrl)
	if err != nil {
		return "", nil, err
	}

	// Определяем MIME-тип файла
	buffer := make([]byte, 512) // Буфер для чтения первых 512 байт
	_, err = file.Read(buffer)
	if err != nil {
		file.Close() // Закрываем файл в случае ошибки чтения
		return "", nil, err
	}

	// Сбрасываем указатель файла в начало для последующего использования
	if _, err := file.Seek(0, 0); err != nil {
		file.Close() // Закрываем файл в случае ошибки
		return "", nil, err
	}

	// Определяем MIME-тип на основе буфера
	mimeType := http.DetectContentType(buffer)

	return mimeType, file, nil
}

func (s *imageService) GetImagesByPropertyId(ctx context.Context, propertyId int64) ([]models.Image, error) {
	images, err := s.imageRepo.GetImagesByPropertyID(ctx, propertyId)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (s *imageService) DeleteImageWithTx(ctx context.Context, imageId int64, tx *sqlx.Tx) error {
	err := s.imageRepo.DeleteWithTx(ctx, imageId, tx)
	if err != nil {
		return err
	}
	return nil
}

func (s *imageService) DeleteImagesByPropertyId(ctx context.Context, propertyId int64, tx *sqlx.Tx) error {
	images, err := s.GetImagesByPropertyId(ctx, propertyId)
	if err != nil {
		return err
	}
	for _, image := range images {
		if err := s.DeleteImageWithTx(ctx, image.Id, tx); err != nil {
			return err
		}
	}
	return nil
}
