package api

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

const destinationFolder = "/static/"

func upload(c echo.Context) error {

	db := c.(*Database)

	user, ok := db.IsAuthorized(c)
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid access rights")
	}

	file, err := c.FormFile("file")

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	src, err := file.Open()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	defer src.Close()

	dst, err := os.Create(destinationFolder + file.Filename)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = db.UpdateUser(user.ID, "", "", "", file.Filename)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func download(c echo.Context) error {
	db := c.(*Database)

	user, ok := db.IsAuthorized(c)
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid access rights")
	}

	filename := user.Avatar

	if filename == "" {
		return c.NoContent(http.StatusOK)
	}

	return c.File(destinationFolder + filename)
}
