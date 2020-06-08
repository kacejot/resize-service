package storage

import (
	"github.com/kacejot/resize-service/pkg/api/graph/model"
	"github.com/kacejot/resize-service/pkg/resize"
)

//go:generate mockgen -destination=mocks/storage_mock.go -package=storage_mock github.com/kacejot/resize-service/pkg/storage Storage

// Storage encapsulates interaction of GraphQL with cloud with images and DB with
// resize records
type Storage interface {
	RecordResizeResult(user string, resizeResult resize.Result) (model.ResizeResult, error)
	GetRecordByID(id string) (original model.Image, err error)
	ListUserRecords(user string) []model.ResizeResult
}
