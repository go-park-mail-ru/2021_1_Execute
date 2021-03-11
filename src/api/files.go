package api

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

const destinationFolder = "../static/"

func saveFile(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destinationFolder + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

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

	err = saveFile(file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = db.UpdateUser(user.ID, "", "", "", destinationFolder+file.Filename)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func download(c echo.Context) error {

	db := c.(*Database)

	_, ok := db.IsAuthorized(c)
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid access rights")
	}

	filename := c.Param("filename")

	if filename == "" {
		return c.NoContent(http.StatusOK)
	}

	return c.File(destinationFolder + filename)
}
