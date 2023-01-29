package testutils

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
	"github.com/stretchr/testify/assert"
)

type RequesterFunc func(r *http.Request) (*http.Response, error)

func (f RequesterFunc) Do(r *http.Request) (*http.Response, error) {
	return f(r)
}

func MockResponse(status int, headers http.Header, body string) *http.Response {
	return &http.Response{
		Status:           strconv.Itoa(status),
		StatusCode:       status,
		Proto:            "HTTP/1.1",
		ProtoMajor:       1,
		ProtoMinor:       1,
		Header:           headers,
		Body:             ioutil.NopCloser(strings.NewReader(body)),
		ContentLength:    0,
		TransferEncoding: []string{},
		Close:            false,
		Uncompressed:     false,
		Trailer:          map[string][]string{},
		Request:          &http.Request{},
		TLS:              &tls.ConnectionState{},
	}
}

func CheckError(t *testing.T, expected apierror.ApiError, actual error) {
	if expected.Code != "" {
		assert.IsType(t, apierror.ApiError{}, actual)
		assert.Equal(t, expected.Code, actual.(apierror.ApiError).Code)
	} else {
		assert.Nil(t, actual)
	}
}
