package interfaces

type ImageRepository interface {
	GetImage(url string) ([]byte, error)
}
