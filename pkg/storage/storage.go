package storage

import "github.com/kacejot/resize-service/pkg/api/graph/model"

// Storage encapsulates interaction of GraphQL with cloud with images and DB with
// resize records
type Storage interface {
	RecordResizeResult(user string, original model.Image, resized model.Image) (model.ResizeResult, error)
	GetRecordByID(id string) (original model.Image, err error)
	ListUserRecords(user string) []model.ResizeResult
}
