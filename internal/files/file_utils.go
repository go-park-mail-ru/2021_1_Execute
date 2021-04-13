package files

import (
	"2021_1_Execute/internal/domain"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const destinationFolder = "/static/"
const localDestinationFolder = "../static/"
const removePrefix = ".."

type fileUtil struct {
}

func NewFileUtil() FileUtil {
	return &fileUtil{}
}

func (fileUtil *fileUtil) GetDestinationFolder() string {
	return destinationFolder
}
func (fileUtil *fileUtil) GetLocalDestinationFolder() string {
	return localDestinationFolder
}

func (fileUtil *fileUtil) DeleteFile(name string) error {
	err := os.Remove(removePrefix + name)
	if err != nil {
		return errors.Wrap(domain.InternalServerError, "Error while deleting file: "+err.Error())
	}
	return nil
}

func (fileUtil *fileUtil) SaveFile(file multipart.File, extension string) (string, error) {

	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(domain.InternalServerError, "Error while create UUID")
	}

	newFilename := uuid.String()
	newFilename += extension

	dst, err := os.Create(localDestinationFolder + newFilename)
	if err != nil {
		return "", errors.Wrap(err, "Error while creating file")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", errors.Wrap(domain.InternalServerError, "Error while saving file")
	}
	return newFilename, nil
}

func (fileUtil *fileUtil) GetExtension(filename string) string {
	splitted := strings.Split(filename, ".")
	if len(splitted) > 1 {
		return "." + splitted[len(splitted)-1]
	}
	return ""
}
