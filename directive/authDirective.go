package directive

import (
	"CP_Discussion/auth"
	"CP_Discussion/database"
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func parseContextClaims(ctx context.Context) (*auth.Claims, error) {
	token, ok := ctx.Value(string("token")).(string)
	if !ok || token == "" {
		return nil, errors.New("no token provided")
	}
	bearer := "Bearer "
	token = token[len(bearer):]
	claims, err := auth.ParseToken(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func AuthDirective(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	_, err := parseContextClaims(ctx)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Access Denied: " + err.Error(),
		}
	}
	return next(ctx)
}

func AdminDirective(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	claims, err := parseContextClaims(ctx)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Access Denied: " + err.Error(),
		}
	}
	member, err := database.DBConnect.FindMemberById(claims.UserID)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Access Denied: member not found",
		}
	}
	if member.IsAdmin {
		return nil, &gqlerror.Error{
			Message: "Access Denied: member is not a admin",
		}
	}
	return next(ctx)
}
