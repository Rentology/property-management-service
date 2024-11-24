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
)

type PropertyService interface {
	Create(ctx context.Context, property *models.Property) (*models.Property, error)
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
