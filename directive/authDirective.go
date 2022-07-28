package directive

import (
	token "CP_Discussion/auth"
	"CP_Discussion/database"
	"CP_Discussion/log"
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func parseContextClaims(ctx context.Context) (*token.Claims, error) {
	_token, ok := ctx.Value("token").(string)
	log.Debug.Println(_token)
	if !ok || _token == "" {
		return nil, errors.New("no token provided")
	}
	bearer := "Bearer "
	_token = _token[len(bearer):]
	claims, err := token.ParseToken(_token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func AuthDirective(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	claims, err := parseContextClaims(ctx)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Access Denied: " + err.Error(),
		}
	}
	ctx = context.WithValue(ctx, "UserID", claims.UserID)
	return next(ctx)
}

func AdminDirective(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	claims, err := parseContextClaims(ctx)
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
