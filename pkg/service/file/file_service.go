package file

import (
	"errors"
	"io/fs"
	"os"

	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
)

type Service struct {
	manager Manager
}

func NewService(manager Manager) Service {
	return Service{
		manager: manager,
	}
}

func NewDefaultService() Service {
	return Service{
		manager: NewDefaultManager(),
	}
}

type Manager interface {
	Stat(name string) (fs.FileInfo, error)
	Mkdir(name string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
	Truncate(name string, size int64) error
	OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error)
}

type manager struct {
	stat      func(name string) (fs.FileInfo, error)
	mkdir     func(name string, perm fs.FileMode) error
	writeFile func(name string, data []byte, perm fs.FileMode) error
	truncate  func(name string, size int64) error
	openFile  func(name string, flag int, perm fs.FileMode) (*os.File, error)
}

func NewDefaultManager() Manager {
	return manager{
		stat:      os.Stat,
		mkdir:     os.Mkdir,
		writeFile: os.WriteFile,
		truncate:  os.Truncate,
		openFile:  os.OpenFile,
	}
}

func NewManager(
	stat func(name string) (fs.FileInfo, error),
	mkdir func(name string, perm fs.FileMode) error,
	writeFile func(name string, data []byte, perm fs.FileMode) error,
	truncate func(name string, size int64) error,
	openFile func(name string, flag int, perm fs.FileMode) (*os.File, error),
) Manager {
	return manager{
		stat:      stat,
		mkdir:     mkdir,
		writeFile: writeFile,
		truncate:  truncate,
		openFile:  openFile,
	}
}

func (m manager) Stat(name string) (fs.FileInfo, error) {
	return m.stat(name)
}

func (m manager) Mkdir(name string, perm fs.FileMode) error {
	return m.mkdir(name, perm)
}

func (m manager) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return m.writeFile(name, data, perm)
}

func (m manager) Truncate(name string, size int64) error {
	return m.truncate(name, size)
}

func (m manager) OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error) {
	return m.openFile(name, flag, perm)
}

func (s Service) CreateDir(path string) error {
	if err := s.manager.Mkdir(path, os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
		return apierror.ApiError{Code: apierror.CreateDirError, InnerCause: err, Values: map[string]any{"path": path}}
	}

	return nil
}

func (s Service) WriteFile(name string, content []byte) error {
	if err := s.manager.WriteFile(name, content, 0666); err != nil {
		return apierror.ApiError{Code: apierror.WriteFileError, InnerCause: err, Values: map[string]any{"name": name}}
	}

	return nil
}

func (s Service) OpenFile(name string, flag int) (*os.File, error) {
	if err := s.manager.Truncate(name, 0); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, apierror.ApiError{Code: apierror.TruncateFileError, InnerCause: err, Values: map[string]any{"name": name}}
	}

	f, err := s.manager.OpenFile(name, flag, 0666)
	if err != nil {
		return nil, apierror.ApiError{Code: apierror.OpenFileError, InnerCause: err, Values: map[string]any{"name": name}}
	}

	return f, nil
}
