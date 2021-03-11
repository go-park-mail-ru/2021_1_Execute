package api

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const destinationFolder = "../static/"

func getExtension(filename string) string {
	splitted := strings.Split(filename, ".")
	if len(splitted) > 1 {
		return "." + splitted[len(splitted)-1]
	}
	return ""
}

func saveFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", errors.Wrap(err, "Error while opening file")
	}
	defer src.Close()

	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "Error while create UUID")
	}

	newFilename := uuid.String()
	newFilename += getExtension(file.Filename)

	dst, err := os.Create(destinationFolder + newFilename)
	if err != nil {
		return "", errors.Wrap(err, "Error while creating file")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", errors.Wrap(err, "Error while saving file")
	}
	return newFilename, nil
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

	filename, err := saveFile(file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = db.UpdateUser(user.ID, "", "", "", destinationFolder+filename)
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
