package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	propertyHttp "property-managment-service/internal/property/delivery/http"
	"property-managment-service/internal/property/repository"
	"property-managment-service/internal/property/service"
	"property-managment-service/pkg/utils"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	propertyRepo := repository.NewPropertyRepository(s.db)
	propertyService := service.NewPropertyService(propertyRepo, s.log)
	propertyHandlers := propertyHttp.NewPropertyHandlers(s.cfg, propertyService, s.log)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true, // разрешает отправку учетных данных
	}))

	v1 := e.Group("/api/v1")
	health := v1.Group("/health")
	propertyGroup := v1.Group("/properties")

	propertyHttp.MapPropertyRoutes(propertyGroup, propertyHandlers, nil)

	health.GET("", func(c echo.Context) error {
		s.log.Info(fmt.Sprintf("Health check RequestID: %s", utils.GetRequestID(c)))
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	return nil

}
