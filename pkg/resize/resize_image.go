package resize

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	rs "github.com/nfnt/resize"
)

// ImageResize implemets Resize
// interface for image resizing
type ImageResize struct{}

// Resize resizes image due to parameters
// png, gif and jpeg are supported
func (ir *ImageResize) Resize(img []byte, width int, height int) (*Result, error) {
	decoded, format, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}

	config, _, err := image.DecodeConfig(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}

	resized := rs.Resize(uint(width), uint(height), decoded, rs.NearestNeighbor)

	buf := new(bytes.Buffer)
	switch format {
	case "png":
		err := png.Encode(buf, resized)
		if err != nil {
			return nil, err
		}
	case "gif":
		err := gif.Encode(buf, resized, &gif.Options{NumColors: 256, Quantizer: nil, Drawer: nil})
		if err != nil {
			return nil, err
		}
	case "jpeg":
		err := jpeg.Encode(buf, resized, &jpeg.Options{Quality: 100})
		if err != nil {
			return nil, err
		}
	}

	return &Result{
		Original: Image{
			Data:   img,
			Width:  config.Width,
			Height: config.Height,
		},
		Resized: Image{
			Data:   buf.Bytes(),
			Width:  width,
			Height: height,
		},
	}, nil
}
