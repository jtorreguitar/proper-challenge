package requesting

import (
	"strconv"
	"sync"

	"github.com/gocolly/colly"
	"github.com/jtorreguitar/proper-challenge/pkg/apierror"
	"github.com/jtorreguitar/proper-challenge/pkg/interfaces"
)

const (
	Class      = ".mu-content-card.mu-card.mu-flush.mu-z1.js-post"
	ResultsDir = "results"
)

type Service struct {
	collector       *colly.Collector
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
	collector *colly.Collector,
	baseUrl string,
	totalImages int,
	fileService interfaces.FileService,
	imageRepo interfaces.ImageRepository,
) Service {
	return Service{
		collector:       collector,
		baseUrl:         baseUrl,
		queue:           &queue{head: &node{action: newScrapeAction(page(baseUrl, 0), collector)}},
		remainingImages: &totalImages,
		totalImages:     totalImages,
		fileService:     fileService,
		imageRepo:       imageRepo,
	}
}

func (s Service) GetImageUrls(maxThreads int) (errorList apierror.ErrorList) {
	if err := s.fileService.CreateDir(ResultsDir); err != nil {
		errorList.List = []apierror.ApiError{wrapErr(err)}
		return errorList
	}

	ch := make(chan error, s.totalImages)
	currentThreads := 0
	currentPage := 0
	var wg sync.WaitGroup
	for s.queue.head != nil {
		if _, ok := s.queue.head.action.(scrapeAction); ok {
			err := s.queue.head.action.a()
			if err != nil {
				return apierror.ErrorList{List: []apierror.ApiError{wrapErr(err)}}
			}

			if *s.remainingImages > 0 {
				currentPage = nextPage(currentPage)
				s.addAction(newScrapeAction(page(s.baseUrl, currentPage), s.collector))
			}
		} else {
			if currentThreads == maxThreads {
				continue
			}

			currentThreads++
			wg.Add(1)
			go func(action action) {
				defer wg.Done()
				if err := action.a(); err != nil {
					ch <- err
				}

				currentThreads--
			}(s.queue.head.action)
		}

		s.queue.head = s.queue.head.next
	}

	wg.Wait()
	close(ch)
	return generateErrorList(ch)
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

func nextPage(current int) int {
	if current == 0 {
		return 2
	}

	return current + 1
}

func generateErrorList(ch chan error) (errorList apierror.ErrorList) {
	for err := range ch {
		errorList.List = append(errorList.List, wrapErr(err))
	}

	return errorList
}
