package http

import (
	"github.com/labstack/echo/v4"
	"property-managment-service/internal/middleware"
)

type PropertyHandlers interface {
	CreateProperty() echo.HandlerFunc
}

func MapPropertyRoutes(userGroup *echo.Group, h PropertyHandlers, mw *middleware.MiddlewareManager) {
	userGroup.POST("", h.CreateProperty())
}
