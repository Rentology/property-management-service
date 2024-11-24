package http

import (
	"github.com/labstack/echo/v4"
	"property-managment-service/internal/middleware"
)

type PropertyHandlers interface {
	CreateProperty() echo.HandlerFunc
	GetProperties() echo.HandlerFunc
}

func MapPropertyRoutes(propertyGroup *echo.Group, h PropertyHandlers, mw *middleware.MiddlewareManager) {
	propertyGroup.POST("", h.CreateProperty())
	propertyGroup.GET("", h.GetProperties())

}
