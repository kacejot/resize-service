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

// RecordResizeResult loads resized images to cloud and creates record about this operation in database
func (is *ImageStorage) RecordResizeResult(ctx context.Context, resizeResult resize.Result) (*model.ResizeResult, error) {
	uploadResult, err := is.cloudStorage.UploadImages(resizeResult)
	if err != nil {
		return nil, err
	}

	return is.recordStorage.StoreRecord(ctx, *uploadResult)
}

// GetRecordByID acquires records of image resize process by ID
func (is *ImageStorage) GetRecordByID(ctx context.Context, id string) (*model.Image, error) {
	result, err := is.recordStorage.FindRecordByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return result.Original, nil
}

// ListUserRecords shows list of resizes that are done by current user
func (is *ImageStorage) ListUserRecords(ctx context.Context) ([]model.ResizeResult, error) {
	return is.recordStorage.FindRecordsForUser(ctx)
}
