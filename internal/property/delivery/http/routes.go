package http

import (
	"github.com/labstack/echo/v4"
	"property-managment-service/internal/middleware"
)

type PropertyHandlers interface {
	CreateProperty() echo.HandlerFunc
	GetProperties() echo.HandlerFunc
	DeleteProperty() echo.HandlerFunc
	UpdateProperty() echo.HandlerFunc
	SavePropertyForm() echo.HandlerFunc
	DeletePropertyForm() echo.HandlerFunc
}

func MapPropertyRoutes(propertyGroup *echo.Group, h PropertyHandlers, mw *middleware.MiddlewareManager) {
	propertyGroup.POST("", h.CreateProperty(), mw.AuthJWTMiddleware())
	propertyGroup.GET("", h.GetProperties())
	propertyGroup.DELETE("/:id", h.DeleteProperty(), mw.AuthJWTMiddleware())
	propertyGroup.PUT("", h.UpdateProperty(), mw.AuthJWTMiddleware())
	propertyGroup.POST("/form", h.SavePropertyForm(), mw.AuthJWTMiddleware())
	propertyGroup.DELETE("/form/:id", h.DeletePropertyForm(), mw.AuthJWTMiddleware())

}
