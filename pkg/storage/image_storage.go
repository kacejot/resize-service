package storage

import (
	"context"

	"github.com/kacejot/resize-service/pkg/api/graph/model"
	"github.com/kacejot/resize-service/pkg/resize"
	"github.com/kacejot/resize-service/pkg/storage/cloud"
	"github.com/kacejot/resize-service/pkg/storage/db"
)

// ImageStorage encapsulates access to cloud and db
// It stores resize results transparently for user
type ImageStorage struct {
	recordStorage db.RecordStorage
	cloudStorage  cloud.Cloud
}

// New creates new instance of ImageStorage
func New() (Storage, error) {
	records, err := db.OpenRecords(db.LoadArangoConfig())
	if err != nil {
		return nil, err
	}

	images := cloud.New(cloud.LoadDropboxConfig())

	return &ImageStorage{
		recordStorage: *records,
		cloudStorage:  *images,
	}, nil
}

// RecordResizeResult loads resized images to cloud and creates record about this operation in database
func (is *ImageStorage) RecordResizeResult(ctx context.Context, resizeResult resize.Result) (*model.ResizeResult, error) {
	uploadResult, err := is.cloudStorage.UploadImages(resizeResult)
	if err != nil {
		return nil, err
	}

	record, err := is.recordStorage.StoreRecord(ctx, *uploadResult)
	if err != nil {
		return nil, err
	}

	return &model.ResizeResult{
		ID: record.ID,
		Original: &model.Image{
			ImageLink: record.Images.Original.ImageLink,
			ExpiresAt: record.Images.Original.ExpiresAt,
			Width:     record.Images.Original.Width,
			Height:    record.Images.Original.Height,
		},
		Resized: &model.Image{
			ImageLink: record.Images.Resized.ImageLink,
			ExpiresAt: record.Images.Resized.ExpiresAt,
			Width:     record.Images.Resized.Width,
			Height:    record.Images.Resized.Height,
		},
	}, nil
}

// GetRecordByID acquires records of image resize process by ID
func (is *ImageStorage) GetRecordByID(ctx context.Context, id string) (*resize.Image, error) {
	result, err := is.recordStorage.FindRecordByID(ctx, id)
	if err != nil {
		return nil, err
	}

	original, err := is.cloudStorage.DownloadImage(result.Images.Original.Path)
	if err != nil {
		return nil, err
	}

	return &resize.Image{
		Data:   original,
		Width:  result.Images.Original.Width,
		Height: result.Images.Original.Height,
	}, nil
}

// ListUserRecords shows list of resizes that are done by current user
func (is *ImageStorage) ListUserRecords(ctx context.Context) ([]*model.ResizeResult, error) {
	records, err := is.recordStorage.FindRecordsForUser(ctx)
	if err != nil {
		return nil, err
	}

	return recordsToResizeResults(records), nil
}

func recordToResizeResult(record *db.Record) *model.ResizeResult {
	if nil == record {
		return nil
	}

	return &model.ResizeResult{
		ID: record.ID,
		Original: &model.Image{
			ImageLink: record.Images.Original.ImageLink,
			ExpiresAt: record.Images.Original.ExpiresAt,
			Width:     record.Images.Original.Width,
			Height:    record.Images.Original.Height,
		},
		Resized: &model.Image{
			ImageLink: record.Images.Resized.ImageLink,
			ExpiresAt: record.Images.Resized.ExpiresAt,
			Width:     record.Images.Resized.Width,
			Height:    record.Images.Resized.Height,
		},
	}
}

func recordsToResizeResults(records []*db.Record) []*model.ResizeResult {
	var result []*model.ResizeResult
	for _, record := range records {
		result = append(result, recordToResizeResult(record))
	}

	return result
}
