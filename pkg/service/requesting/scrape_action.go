package requesting

import (
	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
	"github.com/jtorreguitar/proper-challenge/pkg/interfaces"
)

type scrapeAction struct {
	url     string
	visitor interfaces.VisitorService
}

func newScrapeAction(url string, visitor interfaces.VisitorService) action {
	return scrapeAction{url: url, visitor: visitor}
}

func (action scrapeAction) a() error {
	if err := action.visitor.Visit(action.url); err != nil {
		wrapScrapeErr(err)
	}
	return nil
}

func wrapScrapeErr(err error) error {
	return apierror.ApiError{Code: apierror.ScrapingError, InnerCause: err}
}
