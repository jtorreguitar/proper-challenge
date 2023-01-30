package main

import (
	"flag"
	"fmt"

	"github.com/jtorreguitar/proper-challenge/cmd/api/dependencies"
)

const (
	DefaultImages = 10
)

func main() {
	logger := dependencies.Logger()
	defer logger.Close()

	totalImages := flag.Int("amount", 10, "amount of images to fetch")
	flag.Parse()
	s := dependencies.RequestingService(*totalImages)
	if errorList := s.GetImageUrls(); len(errorList.List) > 0 {
		logger.GenerateErrorReport(errorList)
		fmt.Println("There have been errors. Check logs/error_report.txt for details")
	}
}
