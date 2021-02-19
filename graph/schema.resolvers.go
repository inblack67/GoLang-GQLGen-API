package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/inblack67/GQLGenAPI/graph/generated"
	"github.com/inblack67/GQLGenAPI/graph/model"
)

func (r *queryResolver) Hello(ctx context.Context) (*model.Hello, error) {
	return &model.Hello{
		Reply: "worlds",
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
