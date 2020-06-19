package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"

	"github.com/kacejot/resize-service/pkg/api/graph/generated"
	"github.com/kacejot/resize-service/pkg/api/graph/model"
)

func (r *mutationResolver) Resize(ctx context.Context, image model.ImageInput, width int, height int) (*model.ResizeResult, error) {
	buf, err := base64.StdEncoding.DecodeString(image.Contents)
	if err != nil {
		return nil, err
	}

	result, err := r.imageResize.Resize(buf, width, height)
	if err != nil {
		return nil, err
	}

	response, err := r.imageStorage.RecordResizeResult(ctx, *result)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) ResizeExisting(ctx context.Context, id string, width int, height int) (*model.ResizeResult, error) {
	originalImage, err := r.imageStorage.GetRecordByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result, err := r.imageResize.Resize(originalImage.Data, width, height)
	if err != nil {
		return nil, err
	}

	response, err := r.imageStorage.RecordResizeResult(ctx, *result)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *queryResolver) ListProcessedImages(ctx context.Context) ([]*model.ResizeResult, error) {
	return r.imageStorage.ListUserRecords(ctx)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
