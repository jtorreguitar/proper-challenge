package dependencies

import (
	"net/http"

	"github.com/gocolly/colly"
	"github.com/jtorreguitar/proper-challenge/pkg/interfaces"
	"github.com/jtorreguitar/proper-challenge/pkg/repository/image"
	"github.com/jtorreguitar/proper-challenge/pkg/service/file"
	"github.com/jtorreguitar/proper-challenge/pkg/service/logging"
	"github.com/jtorreguitar/proper-challenge/pkg/service/requesting"
)

func RequestingService(totalImages int) interfaces.RequestingService {
	collector := colly.NewCollector()
	service := requesting.NewService(
		collector,
		"http://icanhas.cheezburger.com",
		totalImages,
		FileService(),
		ImageRepository(),
	)

	collector.OnHTML(requesting.Class, service.GetImageURL)
	return service
}

func FileService() interfaces.FileService {
	return file.NewDefaultService()
}

func ImageRepository() interfaces.ImageRepository {
	return image.NewRepository(HTTPClient())
}

func HTTPClient() *http.Client {
	return http.DefaultClient
}

func Logger() logging.Logger {
	logger, err := logging.NewLogger(file.NewDefaultService(), logging.NewDefaultManager())
	if err != nil {
		panic(err)
	}

	return logger
}
