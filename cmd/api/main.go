package main

import (
	"github.com/gocolly/colly"
	"github.com/jtorreguitar/proper-challenge/pkg/service/requesting"
)

func main() {
	c := colly.NewCollector()
	s := requesting.NewService(c, "http://icanhas.cheezburger.com", 10)
	c.OnHTML(".mu-content-card.mu-card.mu-flush.mu-z1.js-post", s.GetImageUrl)
	s.GetImageUrls()
}
