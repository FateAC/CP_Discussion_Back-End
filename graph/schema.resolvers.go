package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"CP_Discussion/database"
	"CP_Discussion/graph/generated"
	"CP_Discussion/graph/model"
	"context"
)

// CreateMember is the resolver for the createMember field.
func (r *mutationResolver) CreateMember(ctx context.Context, input model.NewMember) (*model.Member, error) {
	return database.DBConnect.InsertMember(input)
}

// LoginCheck is the resolver for the loginCheck field.
func (r *mutationResolver) LoginCheck(ctx context.Context, input model.Login) (*model.Auth, error) {
	return database.DBConnect.LoginCheck(input), nil
}

// Member is the resolver for the member field.
func (r *queryResolver) Member(ctx context.Context, id string) (*model.Member, error) {
	return database.DBConnect.FindMemberById(id)
}

// Members is the resolver for the members field.
func (r *queryResolver) Members(ctx context.Context) ([]*model.Member, error) {
	return database.DBConnect.AllMember()
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
