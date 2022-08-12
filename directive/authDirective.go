package directive

import (
	"CP_Discussion/auth"
	"CP_Discussion/database"
	"CP_Discussion/log"
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func ParseContextClaims(ctx context.Context) (*auth.Claims, error) {
	token, ok := ctx.Value("token").(string)
	log.Debug.Println(token)
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
	claims, err := ParseContextClaims(ctx)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Access Denied: " + err.Error(),
		}
	}
	ctx = context.WithValue(ctx, "UserID", claims.UserID)
	return next(ctx)
}

func AdminDirective(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	claims, err := ParseContextClaims(ctx)
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
