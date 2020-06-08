package resize

//go:generate mockgen -destination=mocks/resize_mock.go -package=resize_mock github.com/kacejot/resize-service/pkg/resize Resize

// Image stores image data encoded in its initial format
type Image struct {
	Data   []byte
	Width  uint
	Height uint
}

// Result is return type of Resize operation
type Result struct {
	Original Image
	Resized  Image
}

// Resize represents the instance that is able to resize an image
type Resize interface {
	Resize(image []byte, width uint, height uint) (Result, error)
}
