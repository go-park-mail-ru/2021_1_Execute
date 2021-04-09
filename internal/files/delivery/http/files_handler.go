package http

import (
	"2021_1_Execute/internal/domain"
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type FilesHandler struct {
	userUC    domain.UserUsecase
	fileUT    domain.FileUtil
	sessionHD domain.SessionHandler
}

func NewFilesHandler(e *echo.Echo, userUsecase domain.UserUsecase,
	fileUtil domain.FileUtil, sessionsHandler domain.SessionHandler) {
	handler := &FilesHandler{
		userUC:    userUsecase,
		fileUT:    fileUtil,
		sessionHD: sessionsHandler,
	}
	e.POST("/api/upload/", handler.AddAvatar)
	e.GET("/static/:filename", handler.Download)
}

func (handler *FilesHandler) AddAvatar(c echo.Context) error {

	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return domain.GetEchoError(err)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return domain.GetEchoError(errors.Wrap(domain.InternalServerError, err.Error()))
	}

	if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image") {
		return domain.GetEchoError(domain.UnsupportedMediaType)
	}

	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		return domain.GetEchoError(errors.Wrap(domain.InternalServerError, err.Error()))
	}
	extension := handler.fileUT.GetExtension(fileHeader.Filename)
	filename, err := handler.fileUT.SaveFile(file, extension)
	if err != nil {
		return domain.GetEchoError(err)
	}
	path := handler.fileUT.GetDestinationFolder() + filename

	ctx := context.Background()
	err = handler.userUC.UpdateAvatar(ctx, userID, path)
	if err != nil {
		return domain.GetEchoError(err)
	}

	return c.NoContent(http.StatusOK)
}

func (handler *FilesHandler) Download(c echo.Context) error {
	_, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return domain.GetEchoError(err)
	}

	filename := c.Param("filename")

	if filename == "" {
		return c.NoContent(http.StatusOK)
	}

	return c.File(handler.fileUT.GetLocalDestinationFolder() + filename)
}
