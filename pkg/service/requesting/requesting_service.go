package requesting

import (
	"strconv"

	"github.com/gocolly/colly"
	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
	"github.com/jtorreguitar/proper-challenge/pkg/interfaces"
)

const (
	Class      = ".mu-content-card.mu-card.mu-flush.mu-z1.js-post"
	ResultsDir = "results"
)

type Service struct {
	visitor         interfaces.VisitorService
	baseUrl         string
	queue           *queue
	remainingImages *int
	totalImages     int
	fileService     interfaces.FileService
	imageRepo       interfaces.ImageRepository
}

type queue struct {
	head *node
	last *node
}

type node struct {
	action action
	next   *node
}

type action interface {
	a() error
}

func NewService(
	visitor interfaces.VisitorService,
	baseUrl string,
	totalImages int,
	fileService interfaces.FileService,
	imageRepo interfaces.ImageRepository,
) Service {
	return Service{
		visitor:         visitor,
		baseUrl:         baseUrl,
		queue:           &queue{head: &node{action: newScrapeAction(page(baseUrl, 0), visitor)}},
		remainingImages: &totalImages,
		totalImages:     totalImages,
		fileService:     fileService,
		imageRepo:       imageRepo,
	}
}

func (s Service) GetImageUrls() (errorList apierror.ErrorList) {
	if err := s.fileService.CreateDir(ResultsDir); err != nil {
		errorList.List = []apierror.ApiError{wrapErr(err)}
		return errorList
	}

	for s.queue.head != nil {
		err := s.queue.head.action.a()

		if err != nil {
			errorList.List = append(errorList.List, wrapErr(err))
		}

		if ae, ok := err.(apierror.ApiError); ok && ae.Code == apierror.ScrapingError {
			return errorList
		}

		s.queue.head = s.queue.head.next
	}

	return errorList
}

func (s Service) GetImageUrl(e *colly.HTMLElement) {
	if *s.remainingImages < 1 {
		return
	}

	url := e.ChildAttr("div a div img", "data-src")
	if url != "" {
		s.addAction(newImageAction(url, ResultsDir, s.totalImages-*s.remainingImages+1, s.imageRepo, s.fileService))
		*s.remainingImages--
	}
}

func (s Service) addAction(a action) {
	newNode := node{action: a}
	if s.queue.last == nil {
		s.queue.head.next = &newNode
		s.queue.last = &newNode
	} else {
		s.queue.last.next = &newNode
		s.queue.last = &newNode
	}
}

func page(url string, page int) string {
	return url + "/page/" + strconv.FormatInt(int64(page), 10)
}

func wrapErr(err error) apierror.ApiError {
	if ae, ok := err.(apierror.ApiError); ok {
		return ae
	}

	return apierror.ApiError{Code: apierror.DefaultError, InnerCause: err}
}
