package middleware

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"property-managment-service/internal/config"
	"property-managment-service/pkg/httpErrors"
)

func (mw *MiddlewareManager) AuthJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Извлекаем токен из куки
			cookieToken, err := c.Cookie("token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.ErrNoCookie))
				}
				return c.JSON(http.StatusInternalServerError, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			token := cookieToken.Value
			if token == "" {
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			// Валидация и разбор токена
			claims, err := mw.validateJWTToken(token, mw.cfg)
			if err != nil {
				mw.log.Error("middleware validateJWTToken", "header JWT", err.Error())
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			// Добавляем данные из токена в контекст
			c.Set("user", claims)

			return next(c)
		}
	}
}

func (mw *MiddlewareManager) validateJWTToken(tokenString string, cfg *config.Config) (map[string]interface{}, error) {
	if tokenString == "" {
		return nil, httpErrors.InvalidJWTToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(cfg.Server.JwtSecretKey)
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, httpErrors.InvalidJWTToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("invalid claims") // todo: Нормальная ошибка
	}

	data := make(map[string]interface{})
	for key, value := range claims {
		data[key] = value
	}

	return data, nil

}
