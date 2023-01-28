package requesting

import (
	"github.com/gocolly/colly"
	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
)

type scrapeAction struct {
	url       string
	collector *colly.Collector
}

func newScrapeAction(url string, collector *colly.Collector) action {
	return scrapeAction{url: url, collector: collector}
}

func (action scrapeAction) a() error {
	return wrapScrapeErr(action.collector.Visit(action.url))
}

func wrapScrapeErr(err error) error {
	return apierror.ApiError{Code: apierror.ScrapingError, InnerCause: err}
}
