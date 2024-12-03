package http

import (
	"github.com/labstack/echo/v4"
	"property-managment-service/internal/middleware"
)

type PropertyDetailsHandlers interface {
	CreatePropertyDetails() echo.HandlerFunc
	GetPropertyDetailsById() echo.HandlerFunc
}

func MapPropertyDetailsRoutes(propertyGroup *echo.Group, h PropertyDetailsHandlers, mw *middleware.MiddlewareManager) {
	propertyGroup.POST("", h.CreatePropertyDetails())
	propertyGroup.GET("/:id", h.GetPropertyDetailsById())
}
