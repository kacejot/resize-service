package resize

//go:generate mockgen -destination=mocks/resize_mock.go -package=resize_mock github.com/kacejot/resize-service/pkg/resize Resize

// Image stores image data encoded in its initial format
type Image struct {
	Data   []byte
	Width  int
	Height int
}

// Result is return type of Resize operation
type Result struct {
	Original Image
	Resized  Image
}

// Resize represents the instance that is able to resize an image
type Resize interface {
	Resize(image []byte, width int, height int) (*Result, error)
}
