package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
	"github.com/jtorreguitar/proper-challenge/pkg/interfaces"
)

type Logger struct {
	manager Manager
	output  *os.File
}

type Manager interface {
	SetOutput(w io.Writer)
	Print(v ...any)
}

type manager struct {
	setOutput func(w io.Writer)
	print     func(v ...any)
}

func (m manager) SetOutput(w io.Writer) {
	m.setOutput(w)
}

func (m manager) Print(v ...any) {
	m.print(v)
}

func NewLogger(fileService interfaces.FileService, manager Manager) (Logger, error) {
	if err := fileService.CreateDir("logs"); err != nil {
		return Logger{}, err
	}

	f, err := fileService.OpenFile("logs/error_report.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY)
	if err != nil {
		return Logger{}, err
	}

	manager.SetOutput(f)
	return Logger{manager: manager, output: f}, nil
}

func NewDefaultManager() Manager {
	return manager{
		setOutput: log.SetOutput,
		print:     log.Print,
	}
}

func NewManager(setOutput func(w io.Writer), print func(v ...any)) Manager {
	return manager{
		setOutput: setOutput,
		print:     print,
	}
}

func (logger Logger) Close() {
	logger.output.Close()
}

func (logger Logger) GenerateErrorReport(errorList apierror.ErrorList) {
	for _, err := range errorList.List {
		logger.manager.Print(generateErroString(err))
	}
}

func generateErroString(err apierror.ApiError) string {
	var builder strings.Builder
	builder.WriteString("\n" + err.Error() + "\n")
	for k, v := range err.Values {
		builder.WriteString(fmt.Sprintf(k+": %+v\n", v))
	}
	builder.WriteString("\n")
	return builder.String()
}
