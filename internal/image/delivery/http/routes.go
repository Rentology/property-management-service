package http

import (
	"github.com/labstack/echo/v4"
	"property-managment-service/internal/middleware"
)

type ImageHandlers interface {
	UploadImage() echo.HandlerFunc
	GetImage() echo.HandlerFunc
	GetImageByPropertyId() echo.HandlerFunc
}

func MapImageRoutes(imageGroup *echo.Group, h ImageHandlers, mw *middleware.MiddlewareManager) {
	imageGroup.POST("", h.UploadImage())
	imageGroup.GET("/:id", h.GetImage())
	imageGroup.GET("", h.GetImageByPropertyId())
}
