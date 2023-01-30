package requesting_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gocolly/colly"
	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
	"github.com/jtorreguitar/proper-challenge/pkg/service/requesting"
	"github.com/jtorreguitar/proper-challenge/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/page/0", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
						<html>
							<div class="mu-content-card mu-card mu-flush mu-z1 js-post">
								<div>
									<a>
										<div>
									   		<img data-src="https://i.chzbgr.com/full/9728348928/h80C82339/maybe-shes-born-with-it-maybe-its-meowbelline"> 
										</div>
									</a>
								</div>
							</div>
						</html>`,
		),
		)
	})

	return httptest.NewServer(mux)
}

func Test_GetImageUrls(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	totalImages := 10
	tests := []struct {
		name        string
		expectedErr apierror.ApiError
		fileService fileServiceMock
		imageRepo   imageRepoMock
	}{
		{
			name: "success",
			fileService: fileServiceMock{
				createDir: func(path string) error { return nil },
				writeFile: func(name string, content []byte) error { return nil },
			},
			imageRepo: func(url string) ([]byte, error) { return []byte{}, nil },
		},
		{
			name:        "fail (create dir)",
			expectedErr: apierror.ApiError{Code: apierror.CreateDirError},
			fileService: fileServiceMock{
				createDir: func(path string) error { return apierror.ApiError{Code: apierror.CreateDirError} },
				writeFile: func(name string, content []byte) error { return nil },
			},
			imageRepo: func(url string) ([]byte, error) { return []byte{}, nil },
		},
		{
			name:        "fail (write file)",
			expectedErr: apierror.ApiError{Code: apierror.WriteFileError},
			fileService: fileServiceMock{
				createDir: func(path string) error { return nil },
				writeFile: func(name string, content []byte) error { return apierror.ApiError{Code: apierror.WriteFileError} },
			},
			imageRepo: func(url string) ([]byte, error) { return []byte{}, nil },
		},
		{
			name:        "fail (get image)",
			expectedErr: apierror.ApiError{Code: apierror.GetImageError},
			fileService: fileServiceMock{
				createDir: func(path string) error { return nil },
				writeFile: func(name string, content []byte) error { return nil },
			},
			imageRepo: func(url string) ([]byte, error) { return nil, apierror.ApiError{Code: apierror.GetImageError} },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := colly.NewCollector()
			s := requesting.NewService(c, ts.URL, totalImages, tt.fileService, tt.imageRepo)
			c.OnHTML(requesting.Class, s.GetImageUrl)

			resp := s.GetImageUrls()
			if tt.expectedErr.Code != "" {
				assert.Greater(t, len(resp.List), 0)
				testutils.CheckError(t, tt.expectedErr, resp.List[0])
			} else {
				assert.Equal(t, 0, len(resp.List))
			}
		})
	}
}

type visitorMock struct {
	e        *colly.HTMLElement
	getImage func(e *colly.HTMLElement)
	visit    func(URL string) error
}

func (m visitorMock) Visit(URL string) error {
	return m.visit(URL)
}

type fileServiceMock struct {
	createDir func(path string) error
	writeFile func(name string, content []byte) error
}

func (m fileServiceMock) CreateDir(path string) error {
	return m.createDir(path)
}

func (m fileServiceMock) WriteFile(name string, content []byte) error {
	return m.writeFile(name, content)
}

func (m fileServiceMock) OpenFile(name string, flag int) (*os.File, error) {
	return nil, nil
}

type imageRepoMock func(url string) ([]byte, error)

func (f imageRepoMock) GetImage(url string) ([]byte, error) {
	return f(url)
}
