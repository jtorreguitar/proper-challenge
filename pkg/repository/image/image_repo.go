package image

import (
	"io"
	"net/http"

	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
	"github.com/jtorreguitar/proper-challenge/pkg/utils/requester"
)

type Repository struct {
	requester requester.Requester
}

func NewRepository(requester requester.Requester) Repository {
	return Repository{requester: requester}
}

func (r Repository) GetImage(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, apierror.ApiError{Code: apierror.GenerateImageRequestError, InnerCause: err}
	}

	resp, err := r.requester.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, apierror.ApiError{Code: apierror.GetImageError, InnerCause: err}
	}

	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return data, apierror.ApiError{Code: apierror.ReadImageError, InnerCause: err}
	}

	return data, nil
}
