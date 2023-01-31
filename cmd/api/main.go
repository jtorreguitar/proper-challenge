package main

import (
	"flag"
	"fmt"

	"github.com/jtorreguitar/proper-challenge/cmd/api/dependencies"
)

const (
	DefaultImages  = 10
	DefaultThreads = 5
)

func main() {
	logger := dependencies.Logger()
	defer logger.Close()

	totalImages := flag.Int("amount", 10, "amount of images to fetch")
	threads := flag.Int("threads", 5, "amount of threads to use")
	flag.Parse()
	s := dependencies.RequestingService(getTotalImages(*totalImages))
	if errorList := s.GetImageURLs(getTotalThreads(*threads)); len(errorList.List) > 0 {
		logger.GenerateErrorReport(errorList)
		fmt.Println("There have been errors. Check logs/error_report.txt for details")
	}
}

func getTotalImages(totalImages int) int {
	if totalImages < 1 {
		return DefaultImages
	}

	return totalImages
}

func getTotalThreads(threads int) int {
	if threads < 1 || threads > 5 {
		return DefaultThreads
	}

	return threads
}
