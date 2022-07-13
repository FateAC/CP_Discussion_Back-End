package directive

import (
	"CP_Discussion/auth"
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func AuthDirective(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	claims, ok := ctx.Value(string("auth")).(*auth.Claims)
	if claims == nil || !ok {
		return nil, &gqlerror.Error{
			Message: "Access Denied",
		}
	}
	return next(ctx)
}

func AdminDirective(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	claims, ok := ctx.Value(string("auth")).(*auth.Claims)
	if claims == nil || !ok {
		return nil, &gqlerror.Error{
			Message: "Access Denied",
		}
	}
	// member, err := database.DBConnect.FindMemberById(claims.UserID)
	// log.Debug.Println(member)
	// if err != nil || !member.IsAdmin {
	// 	return nil, &gqlerror.Error{
	// 		Message: "Access Denied",
	// 	}
	// }
	return next(ctx)
}
