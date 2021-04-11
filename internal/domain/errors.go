package domain

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/labstack/echo"
)

type ServerError struct {
	Message string
}

func (e ServerError) Error() string {
	return e.Message
}

var (
	UnauthorizedError    = ServerError{"Unauthorized Error"}
	ServerConflictError  = ServerError{"Something not unique"}
	BadRequestError      = ServerError{"Incorrect data"}
	ServerNotFoundError  = ServerError{"Not found"}
	InternalServerError  = ServerError{"We are not responsible for this"}
	ForbiddenError       = ServerError{"Forbidden"}
	UnsupportedMediaType = ServerError{"Invalid media type"}
)

func GetEchoError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, UnauthorizedError):
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	case errors.Is(err, ServerConflictError):
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	case errors.Is(err, BadRequestError):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case errors.Is(err, ServerNotFoundError):
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case errors.Is(err, ForbiddenError):
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	case errors.Is(err, UnsupportedMediaType):
		return echo.NewHTTPError(http.StatusUnsupportedMediaType, err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

type DatabaseError struct {
	Message string
}

func (e DatabaseError) Error() string {
	return e.Message
}

var (
	DBNotFoundError = DatabaseError{"Not found"}
	DBConflictError = DatabaseError{"Something not unique"}
)

func DBErrorToServerError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, DBNotFoundError):
		return ServerNotFoundError
	case errors.Is(err, DBConflictError):
		return errors.Wrap(ServerConflictError, err.Error())
	default:
		return errors.Wrap(InternalServerError, err.Error())
	}
}
