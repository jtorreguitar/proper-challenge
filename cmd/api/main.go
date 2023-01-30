package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jtorreguitar/proper-challenge/cmd/api/dependencies"
)

const (
	DefaultImages = 10
)

func main() {
	logger := dependencies.Logger()
	defer logger.Close()

	s := dependencies.RequestingService(getTotalImages())
	if errorList := s.GetImageUrls(); len(errorList.List) > 0 {
		logger.GenerateErrorReport(errorList)
		fmt.Println("There have been errors. Check logs/error_report.txt for details")
	}
}

func getTotalImages() int {
	args := os.Args[1:]
	if len(args) > 0 {
		totalImages, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		return totalImages
	}

	return DefaultImages
}
