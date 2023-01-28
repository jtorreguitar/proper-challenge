package file

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
)

func CreateDir(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return apierror.ApiError{Code: apierror.CreateDirError, InnerCause: err}
		}
	}
	return nil
}

func WriteFile(name string, content []byte) error {
	if err := ioutil.WriteFile(name, content, 0666); err != nil {
		return apierror.ApiError{Code: apierror.WriteFileError, InnerCause: err}
	}

	return nil
}
