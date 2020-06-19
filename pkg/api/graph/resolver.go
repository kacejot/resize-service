package graph

import (
	"github.com/kacejot/resize-service/pkg/resize"
	"github.com/kacejot/resize-service/pkg/storage"
)

// Resolver stores context required for query and mutation resolvers
type Resolver struct {
	imageResize  resize.Resize
	imageStorage storage.Storage
}

// NewResolver returns new instance of Resolver
func NewResolver() (*Resolver, error) {
	s, err := storage.New()
	if err != nil {
		return nil, err
	}

	return &Resolver{
		imageResize:  &resize.ImageResize{},
		imageStorage: s,
	}, nil
}
