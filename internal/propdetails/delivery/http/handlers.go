package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"property-managment-service/internal/models"
	"property-managment-service/pkg/httpErrors"
	"property-managment-service/pkg/utils"
)

type PropertyDetailsService interface {
	Create(ctx context.Context, details *models.PropertyDetails) (*models.PropertyDetails, error)
}

type propertyDetailsHandlers struct {
	propertyDetailsService PropertyDetailsService
	log                    *slog.Logger
}

func NewPropertyDetailsHandlers(service PropertyDetailsService, log *slog.Logger) PropertyDetailsHandlers {
	return &propertyDetailsHandlers{propertyDetailsService: service, log: log}
}

func (h *propertyDetailsHandlers) CreatePropertyDetails() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		requestID := utils.GetRequestID(c)
		h.log.Info("Handling Create", slog.String("request_id", requestID))
		details := &models.PropertyDetails{}
		if err := utils.ReadRequest(c, details); err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		details, err := h.propertyDetailsService.Create(ctx, details)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, details)
	}
}
