package interfaces

import (
	"os"

	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
)

type FileService interface {
	CreateDir(path string) error
	WriteFile(name string, content []byte) error
	OpenFile(name string, flag int) (*os.File, error)
}

type RequestingService interface {
	GetImageUrls() (errorList apierror.ErrorList)
}
