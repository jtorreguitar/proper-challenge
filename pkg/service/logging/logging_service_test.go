package logging_test

import (
	"io"
	"os"
	"testing"

	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
	"github.com/jtorreguitar/proper-challenge/pkg/service/logging"
	"github.com/jtorreguitar/proper-challenge/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func Test_NewLogger(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr apierror.ApiError
		fileService fileServiceMock
		setOutput   func(w io.Writer)
	}{
		{
			name: "success",
			fileService: fileServiceMock{
				createDir: func(path string) error {
					return nil
				},
				openFile: func(name string, flag int) (*os.File, error) {
					return nil, nil
				},
			},
			setOutput: func(w io.Writer) {},
		},
		{
			name:        "fail (create dir)",
			expectedErr: apierror.ApiError{Code: apierror.CreateDirError},
			fileService: fileServiceMock{
				createDir: func(path string) error {
					return apierror.ApiError{Code: apierror.CreateDirError}
				},
			},
		},
		{
			name:        "fail (open file)",
			expectedErr: apierror.ApiError{Code: apierror.OpenFileError},
			fileService: fileServiceMock{
				createDir: func(path string) error {
					return nil
				},
				openFile: func(name string, flag int) (*os.File, error) {
					return nil, apierror.ApiError{Code: apierror.OpenFileError}
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := logging.NewManager(tt.setOutput, nil)

			_, err := logging.NewLogger(tt.fileService, manager)

			testutils.CheckError(t, tt.expectedErr, err)
		})
	}
}

func Test_GenerateErrorReport(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		err      apierror.ApiError
	}{
		{
			name:     "success",
			expected: "\n" + apierror.DefaultError + "\n\n",
			err:      apierror.ApiError{Code: apierror.DefaultError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := ""
			managerMock := loggerManager{output: &output}
			manager := logging.NewManager(func(w io.Writer) {}, managerMock.print)
			logger, _ := logging.NewLogger(
				fileServiceMock{createDir: func(path string) error { return nil }, openFile: func(name string, flag int) (*os.File, error) { return nil, nil }},
				manager,
			)

			logger.GenerateErrorReport(apierror.ErrorList{List: []apierror.ApiError{tt.err}})
			assert.Equal(t, tt.expected, *managerMock.output)
		})
	}
}

type fileServiceMock struct {
	createDir func(path string) error
	openFile  func(name string, flag int) (*os.File, error)
}

func (mock fileServiceMock) CreateDir(path string) error {
	return mock.createDir(path)
}

func (mock fileServiceMock) OpenFile(name string, flag int) (*os.File, error) {
	return mock.openFile(name, flag)
}

func (mock fileServiceMock) WriteFile(name string, content []byte) error {
	return nil
}

type loggerManager struct {
	output *string
}

func (m loggerManager) print(v ...any) {
	if len(v) > 0 {
		*m.output = v[0].([]interface{})[0].(string)
	}
}
