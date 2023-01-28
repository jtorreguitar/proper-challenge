package logging

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
	"github.com/jtorreguitar/proper-challenge/pkg/service/file"
)

type Logger struct {
	output *os.File
}

func Init() (Logger, error) {
	if err := file.CreateDir("logs"); err != nil {
		return Logger{}, err
	}

	f, err := file.OpenFile("logs/error_report.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY)
	if err != nil {
		return Logger{}, err
	}

	log.SetOutput(f)
	return Logger{output: f}, nil
}

func (logger Logger) Close() {
	logger.output.Close()
}

func (logger Logger) GenerateErrorReport(errorList apierror.ErrorList) {
	for _, err := range errorList.List {
		log.Print(generateErroString(err))
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
