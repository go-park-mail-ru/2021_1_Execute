package domain

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
)

type ServerError struct {
	Message string
}

func (e ServerError) Error() string {
	return e.Message
}

var (
	UnauthorizedError   = ServerError{"Unauthorized Error"}
	ConflictError       = ServerError{"Something not unique"}
	BadRequestError     = ServerError{"Incorrect data"}
	NotFoundError       = ServerError{"Not found"}
	InternalServerError = ServerError{"We are not responsible for this"}
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
