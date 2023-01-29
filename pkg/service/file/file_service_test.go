package file_test

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
	"github.com/jtorreguitar/proper-challenge/pkg/service/file"
	"github.com/jtorreguitar/proper-challenge/pkg/testutils"
)

func Test_CreateDir(t *testing.T) {
	dir := "test"
	tests := []struct {
		name        string
		expectedErr apierror.ApiError
		stat        func(name string) (fs.FileInfo, error)
		mkdir       func(name string, perm fs.FileMode) error
	}{
		{
			name: "success",
			stat: func(name string) (fs.FileInfo, error) {
				return nil, nil
			},
			mkdir: func(name string, perm fs.FileMode) error {
				return nil
			},
		},
		{
			name: "success (dir not exists)",
			stat: func(name string) (fs.FileInfo, error) {
				return nil, os.ErrNotExist
			},
			mkdir: func(name string, perm fs.FileMode) error {
				return nil
			},
		},
		{
			name:        "fail (stat)",
			expectedErr: apierror.ApiError{Code: apierror.StatDirError},
			stat: func(name string) (fs.FileInfo, error) {
				return nil, errors.New("hardcoded")
			},
		},
		{
			name:        "fail (mkdir)",
			expectedErr: apierror.ApiError{Code: apierror.CreateDirError},
			stat: func(name string) (fs.FileInfo, error) {
				return nil, nil
			},
			mkdir: func(name string, perm fs.FileMode) error {
				return errors.New("hardcoded")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := file.NewService(file.NewManager(
				tt.stat,
				tt.mkdir,
				nil,
				nil,
				nil,
			))

			err := s.CreateDir(dir)
			testutils.CheckError(t, tt.expectedErr, err)
		})
	}
}

func Test_WriteFile(t *testing.T) {
	fileName := "test"
	content := []byte{}

	tests := []struct {
		name        string
		expectedErr apierror.ApiError
		writeFile   func(name string, data []byte, perm fs.FileMode) error
	}{
		{
			name: "success",
			writeFile: func(name string, data []byte, perm fs.FileMode) error {
				return nil
			},
		},
		{
			name:        "fail",
			expectedErr: apierror.ApiError{Code: apierror.WriteFileError},
			writeFile: func(name string, data []byte, perm fs.FileMode) error {
				return errors.New("hardcoded")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := file.NewService(file.NewManager(
				nil,
				nil,
				tt.writeFile,
				nil,
				nil,
			))

			err := s.WriteFile(fileName, content)
			testutils.CheckError(t, tt.expectedErr, err)
		})
	}
}

func Test_OpenFile(t *testing.T) {
	fileName := "test"
	flag := 0
	tests := []struct {
		name        string
		expectedErr apierror.ApiError
		truncate    func(name string, size int64) error
		openFile    func(name string, flag int, perm fs.FileMode) (*os.File, error)
	}{
		{
			name: "success",
			truncate: func(name string, size int64) error {
				return nil
			},
			openFile: func(name string, flag int, perm fs.FileMode) (*os.File, error) {
				return nil, nil
			},
		},
		{
			name:        "fail (truncate)",
			expectedErr: apierror.ApiError{Code: apierror.TruncateFileError},
			truncate: func(name string, size int64) error {
				return errors.New("hardcoded")
			},
		},
		{
			name:        "fail (open)",
			expectedErr: apierror.ApiError{Code: apierror.OpenFileError},
			truncate: func(name string, size int64) error {
				return nil
			},
			openFile: func(name string, flag int, perm fs.FileMode) (*os.File, error) {
				return nil, errors.New("hardcoded")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := file.NewService(file.NewManager(
				nil,
				nil,
				nil,
				tt.truncate,
				tt.openFile,
			))

			_, err := s.OpenFile(fileName, flag)
			testutils.CheckError(t, tt.expectedErr, err)
		})
	}
}
