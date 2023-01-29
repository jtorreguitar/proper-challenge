package requesting

import (
	"strconv"

	"github.com/jtorreguitar/proper-challenge/pkg/interfaces"
)

type imageAction struct {
	url         string
	resultsDir  string
	number      int
	imageRepo   interfaces.ImageRepository
	fileService interfaces.FileService
}

func newImageAction(
	url,
	resultsDir string,
	number int,
	imageRepo interfaces.ImageRepository,
	fileService interfaces.FileService,
) action {
	return imageAction{
		url:         url,
		resultsDir:  resultsDir,
		number:      number,
		imageRepo:   imageRepo,
		fileService: fileService,
	}
}

func (action imageAction) a() error {
	image, err := action.imageRepo.GetImage(action.url)
	if err != nil {
		return err
	}

	return action.fileService.WriteFile(action.resultsDir+"/"+strconv.FormatInt(int64(action.number), 10)+".jpg", image)
}
