package graph

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kacejot/resize-service/pkg/api/graph/model"
	"github.com/kacejot/resize-service/pkg/resize"
	resize_mock "github.com/kacejot/resize-service/pkg/resize/mocks"
	storage_mock "github.com/kacejot/resize-service/pkg/storage/mocks"
	"github.com/kacejot/resize-service/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const validPng = ` 	
    89 50 4e 47 0d 0a 1a 0a 00 00 00 0d 49 48 44 52
    00 00 00 01 00 00 00 01 08 02 00 00 00 90 77 53
    de 00 00 00 01 73 52 47 42 00 ae ce 1c e9 00 00
    00 04 67 41 4d 41 00 00 b1 8f 0b fc 61 05 00 00
    00 09 70 48 59 73 00 00 0e c3 00 00 0e c3 01 c7
    6f a8 64 00 00 00 0c 49 44 41 54 18 57 63 78 2b
    a3 02 00 03 27 01 2e 15 6b be e9 00 00 00 00 49
    45 4e 44 ae 42 60 82`

func TestResizeValidImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*10)
	ctx = context.WithValue(ctx, UserContextKey, "sample_user")
	defer cancelFn()

	mockResize := resize_mock.NewMockResize(mockCtrl)
	mockStorage := storage_mock.NewMockStorage(mockCtrl)

	imageBuffer, err := hex.DecodeString(utils.RemoveWhitepaces(validPng))
	assert.Nil(t, err)

	resizeResult := resize.Result{
		Original: resize.Image{
			Data:   imageBuffer,
			Width:  1,
			Height: 1,
		},
		Resized: resize.Image{
			Data:   imageBuffer,
			Width:  480,
			Height: 320,
		},
	}

	resizeCall := mockResize.EXPECT().
		Resize(imageBuffer, 480, 320).
		Times(1).
		Return(resizeResult, nil)

	recordResult := model.ResizeResult{
		ID: "1",
		Original: &model.Image{
			ImageLink: "http://link.to.original.image.com",
			Width:     1,
			Height:    1,
		},
		Resized: &model.Image{
			ImageLink: "http://link.to.resized.image.com",
			Width:     480,
			Height:    320,
		},
	}

	mockStorage.EXPECT().
		RecordResizeResult("sample_user", resizeResult).
		Times(1).
		After(resizeCall).
		Return(recordResult, nil)

	resolver := &Resolver{
		imageResize:  mockResize,
		imageStorage: mockStorage,
	}

	response, err := resolver.Mutation().Resize(
		ctx,
		model.ImageInput{
			Filename: "sample.png",
			Contents: string(imageBuffer),
		},
		480,
		320)

	assert.Nil(t, err)
	assert.Equal(t, recordResult, *response)
}

func TestResizeInvalidImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*10)
	ctx = context.WithValue(ctx, UserContextKey, "sample_user")
	defer cancelFn()

	mockResize := resize_mock.NewMockResize(mockCtrl)

	imageBuffer, err := hex.DecodeString(utils.RemoveWhitepaces(validPng))
	assert.Nil(t, err)

	// add some garbage in the very start of the image buffer
	imageBuffer = append([]byte{0, 1, 2, 3}, imageBuffer...)

	mockResize.EXPECT().
		Resize(imageBuffer, 480, 320).
		Times(1).
		Return(resize.Result{}, errors.New("image has invalid format"))

	resolver := &Resolver{
		imageResize:  mockResize,
		imageStorage: nil,
	}

	_, err = resolver.Mutation().Resize(
		ctx,
		model.ImageInput{
			Filename: "sample.png",
			Contents: string(imageBuffer),
		},
		480,
		320)

	assert.Error(t, err, "image has invalid format")
}

func TestResizeInvalidImageSize(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*10)
	ctx = context.WithValue(ctx, UserContextKey, "sample_user")
	defer cancelFn()

	imageBuffer, err := hex.DecodeString(utils.RemoveWhitepaces(validPng))
	assert.Nil(t, err)

	mockResize := resize_mock.NewMockResize(mockCtrl)

	width := -10
	height := -5

	mockResize.EXPECT().
		Resize(imageBuffer, width, height).
		Times(1).
		Return(resize.Result{}, errors.New("image has invalid size"))

	resolver := &Resolver{
		imageResize:  mockResize,
		imageStorage: nil,
	}

	_, err = resolver.Mutation().Resize(
		ctx,
		model.ImageInput{
			Filename: "sample.png",
			Contents: string(imageBuffer),
		},
		width,
		height)

	assert.Error(t, err, "image has invalid size")
}

func TestResizeExistingValidImage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*10)
	ctx = context.WithValue(ctx, UserContextKey, "sample_user")
	defer cancelFn()

	// mockResize := resize_mock.NewMockResize(mockCtrl)
	mockStorage := storage_mock.NewMockStorage(mockCtrl)

	mockStorage.EXPECT().
		GetRecordByID("1").
		Times(1).
		Return()
}
