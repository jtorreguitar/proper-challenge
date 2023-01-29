package image_test

import (
	"net/http"
	"testing"

	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
	"github.com/jtorreguitar/proper-challenge/pkg/repository/image"
	"github.com/jtorreguitar/proper-challenge/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_GetImage(t *testing.T) {
	type expectedResponse struct {
		resp []byte
		err  apierror.ApiError
	}

	url := "test"
	tests := []struct {
		name      string
		expected  expectedResponse
		requester testutils.RequesterFunc
	}{
		{
			name:     "success",
			expected: expectedResponse{resp: okResponse()},
			requester: func(r *http.Request) (*http.Response, error) {
				return testutils.MockResponse(http.StatusOK, nil, string(okResponse())), nil
			},
		},
		{
			name:     "fail (do request error)",
			expected: expectedResponse{err: apierror.ApiError{Code: apierror.GetImageError}},
			requester: func(r *http.Request) (*http.Response, error) {
				return testutils.MockResponse(http.StatusBadRequest, nil, ""), nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := image.NewRepository(tt.requester)
			resp, err := repo.GetImage(url)
			if tt.expected.err.Code != "" {
				assert.IsType(t, apierror.ApiError{}, err)
				assert.Equal(t, tt.expected.err.Code, err.(apierror.ApiError).Code)
			} else {
				assert.Equal(t, tt.expected.resp, resp)
			}
		})
	}
}

func okResponse() []byte {
	return []byte("ok")
}
