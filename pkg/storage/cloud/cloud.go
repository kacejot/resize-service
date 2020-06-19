package cloud

import (
	"bytes"
	"io/ioutil"
	"strconv"

	"github.com/google/uuid"
	"github.com/kacejot/resize-service/pkg/resize"
	"github.com/kacejot/resize-service/pkg/utils"
	"github.com/tj/go-dropbox"
)

// Cloud encapsulates dropbox client
type Cloud struct {
	dropboxClient *dropbox.Client
}

// UploadResult contains links to uploaded images and their ID
type UploadResult struct {
	Original UploadedImage
	Resized  UploadedImage
}

// UploadedImage contaons links and meta for uploaded images
type UploadedImage struct {
	Path      string
	ImageLink string
	ExpiresAt string
	Width     int
	Height    int
}

// DropboxConfig stores key and may be updated in future
type DropboxConfig struct {
	Key string
}

// LoadDropboxConfig fills ArangoConfig with cmd args
func LoadDropboxConfig() *DropboxConfig {
	return &DropboxConfig{
		Key: utils.EnvOrDie("DROPBOX_KEY"),
	}
}

// New creates new instance of Dropbox cloud storage
func New(config *DropboxConfig) *Cloud {
	return &Cloud{
		dropboxClient: dropbox.New(dropbox.NewConfig(config.Key)),
	}
}

// UploadImages uploads resized images to Dropbox cloud storage
func (c *Cloud) UploadImages(images resize.Result) (*UploadResult, error) {
	// create unique image names to avoid name collisions
	originalImagePath := "/" + uuid.New().String()
	resizedImagePath := originalImagePath + imageSizeToString(images.Resized) + "." + images.Format

	originalImagePath += "." + images.Format

	originalUploadResult, err := c.uploadImage(originalImagePath, images.Original)
	if err != nil {
		return nil, err
	}

	resizedUploadResult, err := c.uploadImage(resizedImagePath, images.Resized)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		Original: *originalUploadResult,
		Resized:  *resizedUploadResult,
	}, nil
}

// DownloadImage loads original (not resized) image
// to perform another resize in future
func (c *Cloud) DownloadImage(path string) ([]byte, error) {
	downloadResult, err := c.dropboxClient.Files.Download(&dropbox.DownloadInput{
		Path: path,
	})
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(downloadResult.Body)
}

func (c *Cloud) uploadImage(path string, image resize.Image) (*UploadedImage, error) {
	_, err := c.dropboxClient.Files.Upload(&dropbox.UploadInput{
		Path:   path,
		Mode:   dropbox.WriteModeAdd,
		Reader: bytes.NewReader(image.Data),
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.dropboxClient.Sharing.CreateSharedLink(&dropbox.CreateSharedLinkInput{
		Path: path,
	})
	if err != nil {
		return nil, err
	}

	return &UploadedImage{
		Path:      path,
		ImageLink: resp.URL,
		ExpiresAt: resp.Expires.String(),
		Width:     image.Width,
		Height:    image.Height,
	}, nil
}

func imageSizeToString(image resize.Image) string {
	return strconv.Itoa(image.Width) + "x" + strconv.Itoa(image.Height)
}
