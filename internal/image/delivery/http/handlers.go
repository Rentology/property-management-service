package http

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"property-managment-service/internal/config"
	"property-managment-service/internal/models"
	"property-managment-service/pkg/httpErrors"
	"property-managment-service/pkg/utils"
	"strconv"
)

type ImageService interface {
	UploadImage(ctx context.Context, file *multipart.FileHeader, propertyId int64) error
	GetImage(ctx context.Context, id int64) (string, *os.File, error)
	GetImagesByPropertyId(ctx context.Context, propertyId int64) ([]models.Image, error)
	UploadImageFromBase64(ctx context.Context, base64Image string, propertyId int64, tx *sqlx.Tx) error
	UploadImagesFromBase64(ctx context.Context, base64Images []string, propertyId int64, tx *sqlx.Tx) error
}

type imageHandlers struct {
	cfg          *config.Config
	log          *slog.Logger
	imageService ImageService
}

func NewImageHandlers(cfg *config.Config, service ImageService, log *slog.Logger) ImageHandlers {
	return &imageHandlers{cfg: cfg, log: log, imageService: service}
}

func (h *imageHandlers) UploadImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		propertyIdStr := c.FormValue("propertyId")
		if propertyIdStr == "" {
			return c.JSON(http.StatusBadRequest, "Отсутствует propertyId")
		}

		propertyId, err := strconv.ParseInt(propertyIdStr, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid propertyId")
		}

		file, err := c.FormFile("image")
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Не удалось получить файл")
		}
		err = h.imageService.UploadImage(ctx, file, propertyId)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, "internal error")
		}
		return c.JSON(http.StatusOK, "success")
	}
}

func (h *imageHandlers) GetImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID")
		}
		mimeType, file, err := h.imageService.GetImage(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.String(http.StatusNotFound, "Image not found")
			}
			return c.String(http.StatusInternalServerError, "Error fetching image")
		}
		defer file.Close()
		return c.Stream(http.StatusOK, mimeType, file)
	}
}

func (h *imageHandlers) GetImageByPropertyId() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		id, err := strconv.ParseInt(c.QueryParam("propertyId"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID")
		}
		images, err := h.imageService.GetImagesByPropertyId(ctx, id)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, images)
	}
}
