package api

import (
	"net/http"

	"github.com/labstack/echo"
)

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func GetEchoError(err error) error {
	if err == nil {
		return nil
	}

	switch err.(type) {
	case *UnauthorizedError:
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	case *ConflictError:
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	case *BadRequestError:
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
