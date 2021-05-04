package files

import "mime/multipart"

type FileUtil interface {
	SaveFile(file multipart.File, extension string) (string, error)
	DeleteFile(name string) error
	GetExtension(filename string) string
	GetOriginalFilename(filename string) string
	GetDestinationFolder() string
	GetLocalDestinationFolder() string
}
