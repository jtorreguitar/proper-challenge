package requesting

import "github.com/gocolly/colly"

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
	// TODO: wrap err
	return err
}
