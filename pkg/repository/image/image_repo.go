package image

import (
	"io/ioutil"
	"net/http"

	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
)

func GetImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, apierror.ApiError{Code: apierror.GetImageError, InnerCause: err}
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return data, apierror.ApiError{Code: apierror.ReadImageError, InnerCause: err}
}
