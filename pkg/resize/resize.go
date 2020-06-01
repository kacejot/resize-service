package resize

//go:generate mockgen -destination=resize_mock.go -package=resize github.com/kacejot/resize-service/pkg/resize Resize

type Resize interface {
	Resize(image []byte, width uint, height uint) ([]byte, error)
}
