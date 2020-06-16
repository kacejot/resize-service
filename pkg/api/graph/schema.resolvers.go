package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/kacejot/resize-service/pkg/api/graph/generated"
	"github.com/kacejot/resize-service/pkg/api/graph/model"
)

func (r *mutationResolver) Resize(ctx context.Context, image model.ImageInput, width int, height int) (*model.ResizeResult, error) {
	result, err := r.imageResize.Resize([]byte(image.Contents), width, height)
	if err != nil {
		return nil, err
	}

	// user is stored in context by webhook that reads HTTP headers
	user := ctx.Value(UserContextKey).(string)
	response, err := r.imageStorage.RecordResizeResult(user, result)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (r *mutationResolver) ResizeExisting(ctx context.Context, id string, width int, height int) (*model.ResizeResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ListProcessedImages(ctx context.Context) ([]*model.ResizeResult, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
