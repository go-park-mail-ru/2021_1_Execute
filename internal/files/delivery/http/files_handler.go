package http

import (
	"2021_1_Execute/internal/domain"
	"2021_1_Execute/internal/files"
	"2021_1_Execute/internal/session"
	"2021_1_Execute/internal/tasks"
	"2021_1_Execute/internal/users"
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type FilesHandler struct {
	userUC    users.UserUsecase
	fileUT    files.FileUtil
	taskUC    tasks.TaskUsecase
	sessionHD session.SessionHandler
}

func NewFilesHandler(e *echo.Echo, userUsecase users.UserUsecase, taskUsecase tasks.TaskUsecase,
	fileUtil files.FileUtil, sessionsHandler session.SessionHandler) {
	handler := &FilesHandler{
		userUC:    userUsecase,
		fileUT:    fileUtil,
		taskUC:    taskUsecase,
		sessionHD: sessionsHandler,
	}
	e.POST("/api/upload/", handler.AddAvatar)
	e.GET("/static/:filename", handler.Download)
	e.POST("/api/upload/attachment/:taskId/", handler.AddAttachment)
}

func (handler *FilesHandler) AddAvatar(c echo.Context) error {

	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return errors.Wrap(domain.InternalServerError, err.Error())
	}

	if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image") {
		return domain.UnsupportedMediaType
	}

	file, err := fileHeader.Open()
	if err != nil {
		return errors.Wrap(domain.InternalServerError, err.Error())
	}
	defer file.Close()
	extension := handler.fileUT.GetExtension(fileHeader.Filename)
	filename, err := handler.fileUT.SaveFile(file, extension)
	if err != nil {
		return err
	}
	path := handler.fileUT.GetDestinationFolder() + filename

	ctx := context.Background()
	err = handler.userUC.UpdateAvatar(ctx, userID, path)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (handler *FilesHandler) Download(c echo.Context) error {
	_, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	filename := c.Param("filename")

	if filename == "" {
		return c.NoContent(http.StatusOK)
	}

	return c.File(handler.fileUT.GetLocalDestinationFolder() + filename)
}

func (handler *FilesHandler) AddAttachment(c echo.Context) error {
	userID, err := handler.sessionHD.IsAuthorized(c)
	if err != nil {
		return err
	}

	taskID, err := strconv.Atoi(c.Param("taskId"))
	if err != nil || taskID < 0 {
		return domain.IDFormatError
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return errors.Wrap(domain.InternalServerError, err.Error())
	}

	file, err := fileHeader.Open()
	if err != nil {
		return errors.Wrap(domain.InternalServerError, err.Error())
	}
	defer file.Close()

	extension := handler.fileUT.GetExtension(fileHeader.Filename)
	originalFilename := handler.fileUT.GetOriginalFilename(fileHeader.Filename)

	filename, err := handler.fileUT.SaveFile(file, extension)
	if err != nil {
		return err
	}

	path := handler.fileUT.GetDestinationFolder() + filename

	_, err = handler.taskUC.AddAttachment(
		context.Background(),
		taskID,
		tasks.Attachment{
			Name: originalFilename,
			Path: path,
		},
		userID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
