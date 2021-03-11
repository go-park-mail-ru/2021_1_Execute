package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
)

var (
	UnauthorizedError = errors.New("")
	ConflictError     = errors.New("")
	BadRequestError   = errors.New("")
	NotFoundError     = errors.New("")
)

func GetEchoError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, UnauthorizedError):
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	case errors.Is(err, ConflictError):
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	case errors.Is(err, BadRequestError):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case errors.Is(err, NotFoundError):
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
