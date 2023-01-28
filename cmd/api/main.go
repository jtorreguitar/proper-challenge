package main

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/jtorreguitar/proper-challenge/pkg/service/logging"
	"github.com/jtorreguitar/proper-challenge/pkg/service/requesting"
)

func main() {
	logger, err := logging.Init()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer logger.Close()

	c := colly.NewCollector()
	s := requesting.NewService(c, "http://icanhas.cheezburger.com", 10)
	c.OnHTML(".mu-content-card.mu-card.mu-flush.mu-z1.js-post", s.GetImageUrl)
	if errorList := s.GetImageUrls(); len(errorList.List) > 0 {
		logger.GenerateErrorReport(errorList)
		fmt.Println("There have been errors. Check logs/error_report.txt for details")
	}
}
