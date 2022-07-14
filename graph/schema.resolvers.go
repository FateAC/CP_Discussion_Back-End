package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"CP_Discussion/auth"
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

// AddMemberCourse is the resolver for the addMemberCourse field.
func (r *mutationResolver) AddMemberCourse(ctx context.Context, id string, course model.NewCourse) (*model.Member, error) {
	return database.DBConnect.AddMemberCourse(id, course)
}

// RemoveMemberCourse is the resolver for the removeMemberCourse field.
func (r *mutationResolver) RemoveMemberCourse(ctx context.Context, id string, course model.NewCourse) (*model.Member, error) {
	return database.DBConnect.RemoveMemberCourse(id, course)
}

// Member is the resolver for the member field.
func (r *queryResolver) Member(ctx context.Context, id string) (*model.Member, error) {
	return database.DBConnect.FindMemberById(id)
}

// Members is the resolver for the members field.
func (r *queryResolver) Members(ctx context.Context) ([]*model.Member, error) {
	return database.DBConnect.AllMember()
}

// IsAdmin is the resolver for the isAdmin field.
func (r *queryResolver) IsAdmin(ctx context.Context) (bool, error) {
	claims, _ := ctx.Value(string("auth")).(*auth.Claims)
	return database.DBConnect.MemberIsAdmin(claims.UserID), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
