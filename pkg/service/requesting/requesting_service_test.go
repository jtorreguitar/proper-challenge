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

	mux.HandleFunc("/page/2", func(w http.ResponseWriter, r *http.Request) {
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

	tests := []struct {
		name        string
		totalImages int
		expectedErr apierror.ApiError
		fileService fileServiceMock
		imageRepo   imageRepoMock
	}{
		{
			name:        "success",
			totalImages: 1,
			fileService: fileServiceMock{
				createDir: func(path string) error { return nil },
				writeFile: func(name string, content []byte) error { return nil },
			},
			imageRepo: func(url string) ([]byte, error) { return []byte{}, nil },
		},
		{
			name:        "fail (create dir)",
			totalImages: 1,
			expectedErr: apierror.ApiError{Code: apierror.CreateDirError},
			fileService: fileServiceMock{
				createDir: func(path string) error { return apierror.ApiError{Code: apierror.CreateDirError} },
				writeFile: func(name string, content []byte) error { return nil },
			},
			imageRepo: func(url string) ([]byte, error) { return []byte{}, nil },
		},
		{
			name:        "fail (write file)",
			totalImages: 1,
			expectedErr: apierror.ApiError{Code: apierror.WriteFileError},
			fileService: fileServiceMock{
				createDir: func(path string) error { return nil },
				writeFile: func(name string, content []byte) error { return apierror.ApiError{Code: apierror.WriteFileError} },
			},
			imageRepo: func(url string) ([]byte, error) { return []byte{}, nil },
		},
		{
			name:        "fail (get image)",
			totalImages: 1,
			expectedErr: apierror.ApiError{Code: apierror.GetImageError},
			fileService: fileServiceMock{
				createDir: func(path string) error { return nil },
				writeFile: func(name string, content []byte) error { return nil },
			},
			imageRepo: func(url string) ([]byte, error) { return nil, apierror.ApiError{Code: apierror.GetImageError} },
		},
		{
			name:        "success (more than one page)",
			totalImages: 2,
			fileService: fileServiceMock{
				createDir: func(path string) error { return nil },
				writeFile: func(name string, content []byte) error { return nil },
			},
			imageRepo: func(url string) ([]byte, error) { return []byte{}, nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := colly.NewCollector()
			s := requesting.NewService(c, ts.URL, tt.totalImages, tt.fileService, tt.imageRepo)
			c.OnHTML(requesting.Class, s.GetImageUrl)

			resp := s.GetImageUrls(1)
			if tt.expectedErr.Code != "" {
				assert.Greater(t, len(resp.List), 0)
				testutils.CheckError(t, tt.expectedErr, resp.List[0])
			} else {
				assert.Equal(t, 0, len(resp.List))
			}
		})
	}
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
