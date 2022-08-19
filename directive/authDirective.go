package directive

import (
	"CP_Discussion/auth"
	"CP_Discussion/database"
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func AuthDirective(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	claims, err := auth.ParseContextClaims(ctx)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Access Denied: " + err.Error(),
		}
	}
	ctx = context.WithValue(ctx, "UserID", claims.UserID)
	return next(ctx)
}

func AdminDirective(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	claims, err := auth.ParseContextClaims(ctx)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Access Denied: " + err.Error(),
		}
	}
	ctx = context.WithValue(ctx, "UserID", claims.UserID)
	isAdmin := database.DBConnect.MemberIsAdmin(claims.UserID)
	if !isAdmin {
		return nil, &gqlerror.Error{
			Message: "Access Denied: member is not a admin",
		}
	}
	return next(ctx)
}
