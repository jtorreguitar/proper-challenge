package file

import (
	"errors"
	"os"

	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
)

func CreateDir(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return apierror.ApiError{Code: apierror.CreateDirError, InnerCause: err, Values: map[string]any{"path": path}}
		}
	}
	return nil
}

func WriteFile(name string, content []byte) error {
	if err := os.WriteFile(name, content, 0666); err != nil {
		return apierror.ApiError{Code: apierror.WriteFileError, InnerCause: err, Values: map[string]any{"name": name}}
	}

	return nil
}

func OpenFile(name string, flag int) (*os.File, error) {
	if err := os.Truncate(name, 0); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, apierror.ApiError{Code: apierror.TruncateFileError, InnerCause: err, Values: map[string]any{"name": name}}
	}

	f, err := os.OpenFile(name, flag, 0666)
	if err != nil {
		return nil, apierror.ApiError{Code: apierror.OpenFileError, InnerCause: err, Values: map[string]any{"name": name}}
	}

	return f, nil
}
