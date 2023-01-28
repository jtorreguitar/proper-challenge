package image

import (
	"io"
	"net/http"

	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
)

func GetImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, apierror.ApiError{Code: apierror.GetImageError, InnerCause: err}
	}

	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return data, apierror.ApiError{Code: apierror.ReadImageError, InnerCause: err}
	}

	return data, nil
}
