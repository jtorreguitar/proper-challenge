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
			// TODO: create more verbose error
			return err
		}
	}
	return nil
}

func WriteFile(name string, content []byte) error {
	if err := ioutil.WriteFile(name, content, 0666); err != nil {
		return apierror.ApiError{}
	}

	return nil
}
