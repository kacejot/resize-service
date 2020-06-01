package storage

import "github.com/kacejot/resize-service/pkg/api/graph/model"

//go:generate mockgen -destination=storage_mock.go -package=storage github.com/kacejot/resize-service/pkg/storage Storage

// Storage encapsulates interaction of GraphQL with cloud with images and DB with
// resize records
type Storage interface {
	RecordResizeResult(user string, original model.Image, resized model.Image) (model.ResizeResult, error)
	GetRecordByID(id string) (original model.Image, err error)
	ListUserRecords(user string) []model.ResizeResult
}
