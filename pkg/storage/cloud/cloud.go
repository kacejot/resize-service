package cloud

import (
	"bytes"
	"flag"
	"io/ioutil"
	"strconv"

	"github.com/google/uuid"
	"github.com/kacejot/resize-service/pkg/api/graph/model"
	"github.com/kacejot/resize-service/pkg/resize"
	"github.com/tj/go-dropbox"
)

// Cloud encapsulates dropbox client
type Cloud struct {
	dropboxClient *dropbox.Client
}

// UploadResult contains links to uploaded images and their ID
type UploadResult struct {
	Original *model.Image
	Resized  *model.Image
}

// DropboxConfig stores key and may be updated in future
type DropboxConfig struct {
	Key string
}

// LoadArangoConfig fills ArangoConfig with cmd args
func LoadDropboxConfig() DropboxConfig {
	return DropboxConfig{
		Key: *flag.String("dropbox-key", "", "identifies the storage and the user"),
	}
}

// New creates new instance of Dropbox cloud storage
func New(config DropboxConfig) *Cloud {
	return &Cloud{
		dropboxClient: dropbox.New(&dropbox.Config{
			AccessToken: config.Key,
		}),
	}
}

// UploadImages uploads resized images to Dropbox cloud storage
func (c *Cloud) UploadImages(images resize.Result) (*UploadResult, error) {
	// create unique image names to avoid name collisions
	originalImagePath := uuid.New().String()
	resizedImagePath := originalImagePath + imageSizeToString(images.Resized)

	originalUploadResult, err := c.uploadImage(originalImagePath, images.Original)
	if err != nil {
		return nil, err
	}

	resizedUploadResult, err := c.uploadImage(resizedImagePath, images.Resized)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		Original: originalUploadResult,
		Resized:  resizedUploadResult,
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

func (c *Cloud) uploadImage(path string, image resize.Image) (*model.Image, error) {
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

	return &model.Image{
		ImageLink: resp.URL,
		ExpiresAt: resp.Expires.String(),
		Width:     image.Width,
		Height:    image.Height,
	}, nil
}

func imageSizeToString(image resize.Image) string {
	return strconv.Itoa(image.Width) + "x" + strconv.Itoa(image.Height)
}
