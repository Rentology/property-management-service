package http

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"property-managment-service/internal/config"
	"property-managment-service/internal/models"
	"property-managment-service/internal/models/request"
	"property-managment-service/pkg/httpErrors"
	"property-managment-service/pkg/utils"
	"strconv"
)

type PropertyService interface {
	Create(ctx context.Context, property *models.Property) (*models.Property, error)
	GetById(ctx context.Context, id int64) (*models.Property, error)
	GetByOwnerId(ctx context.Context, id int64) ([]*models.Property, error)
	Delete(ctx context.Context, id int64) (int64, error)
	Update(ctx context.Context, property *models.Property) (*models.Property, error)
	SaveWithTx(ctx context.Context, property *models.Property, tx *sqlx.Tx) error
	DeleteWithTx(ctx context.Context, id int64, tx *sqlx.Tx) error
	GetAll(ctx context.Context) ([]*models.Property, error)
}

type PropertyFormService interface {
	SavePropertyForm(ctx context.Context, form *request.AddPropertyRequest) error
	DeletePropertyForm(ctx context.Context, propertyID int64) error
}

type propertyHandlers struct {
	propertyService     PropertyService
	propertyServiceForm PropertyFormService
	cfg                 *config.Config
	log                 *slog.Logger
}

func NewPropertyHandlers(propertyService PropertyService, propertyServiceForm PropertyFormService, cfg *config.Config, log *slog.Logger) PropertyHandlers {
	return &propertyHandlers{propertyService: propertyService, propertyServiceForm: propertyServiceForm, cfg: cfg, log: log}
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

		userClaims := c.Get("user").(map[string]interface{})
		userIdFromClaims := userClaims["uid"].(float64)

		if (int64(userIdFromClaims)) != property.OwnerId {
			utils.LogResponseError(c, h.log, httpErrors.Unauthorized)
			return c.JSON(http.StatusUnauthorized, httpErrors.Unauthorized)
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
		} else {
			// Если параметры id и ownerId не переданы, возвращаем все записи
			properties, err := h.propertyService.GetAll(ctx)
			if err != nil {
				utils.LogResponseError(c, h.log, err)
				return c.JSON(httpErrors.ErrorResponse(err))
			}
			return c.JSON(http.StatusOK, properties)
		}
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

func (h *propertyHandlers) UpdateProperty() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		requestID := utils.GetRequestID(c)
		h.log.Info("Handling UpdateProperty", slog.String("request_id", requestID))
		property := &models.Property{}
		if err := utils.ReadRequest(c, property); err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		property, err := h.propertyService.Update(ctx, property)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, property)
	}
}

func (h *propertyHandlers) SavePropertyForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		requestID := utils.GetRequestID(c)
		h.log.Info("Handling SavePropertyForm", slog.String("request_id", requestID))
		r := &request.AddPropertyRequest{}

		if err := utils.ReadRequest(c, r); err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		userClaims := c.Get("user").(map[string]interface{})
		userIdFromClaims := userClaims["uid"].(float64)

		r.Property.OwnerId = int64(userIdFromClaims)

		err := h.propertyServiceForm.SavePropertyForm(ctx, r)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, map[string]string{
			"message": "Property form saved successfully",
		})

	}
}

func (h *propertyHandlers) DeletePropertyForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		requestID := utils.GetRequestID(c)
		h.log.Info("Handling DeleteProperty", slog.String("request_id", requestID))
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		userClaims := c.Get("user").(map[string]interface{})
		userIdFromClaims := userClaims["uid"].(float64)

		property, err := h.propertyService.GetById(ctx, id)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if (int64(userIdFromClaims)) != property.OwnerId {
			utils.LogResponseError(c, h.log, httpErrors.Unauthorized)
			return c.JSON(http.StatusUnauthorized, httpErrors.Unauthorized)
		}

		err = h.propertyServiceForm.DeletePropertyForm(ctx, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, map[string]string{
			"message": "Property removed successfully",
		})
	}
}

func (h *propertyHandlers) GetAllProperties() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		requestID := utils.GetRequestID(c)
		h.log.Info("Handling GetAllProperties", slog.String("request_id", requestID))

		// Получаем все записи недвижимости через сервис
		properties, err := h.propertyService.GetAll(ctx)
		if err != nil {
			utils.LogResponseError(c, h.log, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		// Отправляем данные в ответе
		return c.JSON(http.StatusOK, properties)
	}
}
