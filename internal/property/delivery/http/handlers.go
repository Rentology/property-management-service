package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"property-managment-service/internal/config"
	"property-managment-service/internal/models"
	"property-managment-service/pkg/httpErrors"
	"property-managment-service/pkg/utils"
	"strconv"
)

type PropertyService interface {
	Create(ctx context.Context, property *models.Property) (*models.Property, error)
	GetById(ctx context.Context, id int64) (*models.Property, error)
	GetByOwnerId(ctx context.Context, id int64) ([]*models.Property, error)
	Delete(ctx context.Context, id int64) (int64, error)
}

type propertyHandlers struct {
	cfg             *config.Config
	propertyService PropertyService
	log             *slog.Logger
}

func NewPropertyHandlers(cfg *config.Config, propertyService PropertyService, log *slog.Logger) PropertyHandlers {
	return &propertyHandlers{cfg, propertyService, log}
}

func (h *propertyHandlers) CreateProperty() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		requestID := utils.GetRequestID(c)
		h.log.Info("Handling Create", slog.String("request_id", requestID))
		property := &models.Property{}
		if err := utils.ReadRequest(c, property); err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		property, err := h.propertyService.Create(ctx, property)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusCreated, property)
	}
}

func (h *propertyHandlers) GetProperties() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		requestID := utils.GetRequestID(c)
		h.log.Info("Handling GetProperties", slog.String("request_id", requestID))

		// Проверяем, какой параметр был передан
		if idParam := c.QueryParam("id"); idParam != "" {
			// Логика для поиска по property_id
			id, err := strconv.ParseInt(idParam, 10, 64)
			if err != nil {
				utils.LogResponseError(c, h.log, httpErrors.NewBadRequestError("invalid id"))
				return c.JSON(http.StatusBadRequest, httpErrors.NewBadRequestError("invalid id"))
			}

			property, err := h.propertyService.GetById(ctx, id)
			if err != nil {
				utils.LogResponseError(c, h.log, err)
				return c.JSON(httpErrors.ErrorResponse(err))
			}
			return c.JSON(http.StatusOK, property)
		} else if ownerIdParam := c.QueryParam("ownerId"); ownerIdParam != "" {
			// Логика для поиска по owner_id
			ownerId, err := strconv.ParseInt(ownerIdParam, 10, 64)
			if err != nil {
				utils.LogResponseError(c, h.log, httpErrors.NewBadRequestError("invalid ownerId"))
				return c.JSON(http.StatusBadRequest, httpErrors.NewBadRequestError("invalid ownerId"))
			}

			properties, err := h.propertyService.GetByOwnerId(ctx, ownerId)
			if err != nil {
				utils.LogResponseError(c, h.log, err)
				return c.JSON(httpErrors.ErrorResponse(err))
			}
			return c.JSON(http.StatusOK, properties)
		}

		return c.JSON(http.StatusBadRequest, httpErrors.NewBadRequestError("either id or ownerId must be provided"))
	}
}

func (h *propertyHandlers) DeleteProperty() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		requestID := utils.GetRequestID(c)
		h.log.Info("Handling DeleteProperty", slog.String("request_id", requestID))
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		deletedId, err := h.propertyService.Delete(ctx, id)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, deletedId)
	}
}
