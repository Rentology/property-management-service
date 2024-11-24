package utils

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
)

type ReqIDCtxKey struct{}

func GetRequestID(c echo.Context) string {
	return c.Request().Header.Get(echo.HeaderXRequestID)
}

func GetIPAddress(c echo.Context) string {
	return c.Request().RemoteAddr
}

func GetRequestCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), ReqIDCtxKey{}, GetRequestID(c))
}

func ReadRequest(c echo.Context, request interface{}) error {
	if err := c.Bind(request); err != nil {
		return err
	}
	return validate.StructCtx(c.Request().Context(), request)

}

func LogResponseError(ctx echo.Context, log *slog.Logger, err error) {
	log.Error(fmt.Sprintf("ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(ctx),
		GetIPAddress(ctx),
		err,
	))
}
