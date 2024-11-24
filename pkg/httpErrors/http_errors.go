package httpErrors

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	ErrEmailAlreadyExists = "User with given email already exists"
)

var (
	BadRequest          = errors.New("Bad Request")
	Unauthorized        = errors.New("Unauthorized")
	InvalidJWTToken     = errors.New("Invalid JWT token")
	ErrNoCookie         = errors.New("No cookie")
	InvalidJWTClaims    = errors.New("Invalid JWT claims")
	NotFound            = errors.New("Not Found")
	InternalServerError = errors.New("Internal Server Error")
	ExistsEmailError    = errors.New("User with given email already exists")
)

type RestErr interface {
	Status() int
	Error() string
	Causes() interface{}
}

type RestError struct {
	ErrStatus int         `json:"status,omitempty"`
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"-"`
}

func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e RestError) Status() int {
	return e.ErrStatus
}

func (e RestError) Causes() interface{} {
	return e.ErrCauses
}

func NewRestError(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

func NewRestErrorWithMessage(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

func NewInternalServerError(causes interface{}) RestErr {
	result := RestError{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  InternalServerError.Error(),
		ErrCauses: causes,
	}
	return result
}

func NewBadRequestError(causes interface{}) RestErr {
	result := RestError{
		ErrStatus: http.StatusBadRequest,
		ErrError:  BadRequest.Error(),
		ErrCauses: causes,
	}
	return result
}

func NewUnauthorizedError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusUnauthorized,
		ErrError:  Unauthorized.Error(),
		ErrCauses: causes,
	}
}

func ParseErrors(err error) RestErr {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewRestError(http.StatusNotFound, NotFound.Error(), err)
	case strings.Contains(err.Error(), "SQLSTATE"):
		return parseSqlErrors(err)
	default:
		var restErr RestErr
		if errors.As(err, &restErr) {
			return restErr
		}
		return NewInternalServerError(err)
	}
}

func parseSqlErrors(err error) RestErr {
	if strings.Contains(err.Error(), "23505") {
		return NewRestError(http.StatusBadRequest, ExistsEmailError.Error(), err)
	}

	return NewRestError(http.StatusBadRequest, BadRequest.Error(), err)
}

func ErrorResponse(err error) (int, interface{}) {
	return ParseErrors(err).Status(), ParseErrors(err)
}
