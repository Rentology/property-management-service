package service

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	http2 "property-managment-service/internal/image/delivery/http"
	"property-managment-service/internal/models"
)

type ImageRepository interface {
	SaveImage(ctx context.Context, image *models.Image) (*models.Image, error)
	GetImage(ctx context.Context, id int64) (*models.Image, error)
	GetImagesByPropertyID(ctx context.Context, propertyID int64) ([]models.Image, error)
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
