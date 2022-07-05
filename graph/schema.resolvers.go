package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"CP_Discussion/graph/generated"
	"CP_Discussion/graph/model"
	"CP_Discussion/database"
	"context"
)

var db = database.Connect("mongodb://CPDiscussion:94879487@localhost:9487/")

// CreateMovie is the resolver for the createMovie field.
func (r *mutationResolver) CreateMovie(ctx context.Context, input model.NewMovie) (*model.Movie, error) {
	return db.InsertMovieById(input), nil
}

// Movie is the resolver for the movie field.
func (r *queryResolver) Movie(ctx context.Context, id string) (*model.Movie, error) {
	return db.FindMovieById(id), nil
}

// Movies is the resolver for the movies field.
func (r *queryResolver) Movies(ctx context.Context) ([]*model.Movie, error) {
	return db.All(), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
