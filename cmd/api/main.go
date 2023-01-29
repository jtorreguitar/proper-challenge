package main

import (
	"fmt"

	"github.com/jtorreguitar/proper-challenge/cmd/api/dependencies"
)

func main() {
	logger := dependencies.Logger()
	defer logger.Close()

	s := dependencies.RequestingService(10)
	if errorList := s.GetImageUrls(); len(errorList.List) > 0 {
		logger.GenerateErrorReport(errorList)
		fmt.Println("There have been errors. Check logs/error_report.txt for details")
	}
}
