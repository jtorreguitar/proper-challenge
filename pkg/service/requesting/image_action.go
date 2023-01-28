package requesting

import (
	"strconv"

	"github.com/jtorreguitar/proper-challenge/pkg/repository/image"
	"github.com/jtorreguitar/proper-challenge/pkg/service/file"
)

type imageAction struct {
	url    string
	number int
}

func newImageAction(url string, number int) action {
	return imageAction{url: url, number: number}
}

func (action imageAction) a() error {
	image, err := image.GetImage(action.url)
	if err != nil {
		return err
	}

	return file.WriteFile("results/"+strconv.FormatInt(int64(action.number), 10)+".jpg", image)
}
