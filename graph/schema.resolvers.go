package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"CP_Discussion/auth"
	"CP_Discussion/database"
	"CP_Discussion/graph/generated"
	"CP_Discussion/graph/model"
	"CP_Discussion/log"
	"CP_Discussion/mail"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// CreateMember is the resolver for the createMember field.
func (r *mutationResolver) CreateMember(ctx context.Context, input model.NewMember) (*model.Member, error) {
	return database.DBConnect.InsertMember(input)
}

// RemoveMember is the resolver for the removeMember field.
func (r *mutationResolver) RemoveMember(ctx context.Context, id string) (*model.Member, error) {
	return database.DBConnect.DeleteMember(id)
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

// UpdateMemberAvatar is the resolver for the updateMemberAvatar field.
func (r *mutationResolver) UpdateMemberAvatar(ctx context.Context, avatar graphql.Upload) (bool, error) {
	id, ok := ctx.Value(string("UserID")).(string)
	if !ok {
		return false, errors.New("failed to get user id from ctx")
	}
	return database.DBConnect.UpdateMemberAvatar(id, avatar)
}

// RemoveMemberAvatar is the resolver for the removeMemberAvatar field.
func (r *mutationResolver) RemoveMemberAvatar(ctx context.Context, id string) (bool, error) {
	return database.DBConnect.DeleteMemberAvatar(id)
}

// UpdateMemberNickname is the resolver for the updateMemberNickname field.
func (r *mutationResolver) UpdateMemberNickname(ctx context.Context, nickname string) (bool, error) {
	id, ok := ctx.Value(string("UserID")).(string)
	if !ok {
		return false, errors.New("failed to get user id from ctx")
	}
	return database.DBConnect.UpdateMemberNickname(id, nickname)
}

// AddPost is the resolver for the addPost field.
func (r *mutationResolver) AddPost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	return database.DBConnect.InsertPost(input)
}

// RemovePost is the resolver for the removePost field.
func (r *mutationResolver) RemovePost(ctx context.Context, id string) (*model.Post, error) {
	return database.DBConnect.DeletePost(id)
}

// AddPostComment is the resolver for the addPostComment field.
func (r *mutationResolver) AddPostComment(ctx context.Context, id string, newComment model.NewComment) (bool, error) {
	commenterID, ok := ctx.Value(string("UserID")).(string)
	if !ok {
		return false, errors.New("failed to get user id from ctx")
	}
	return database.DBConnect.AddPostComment(id, commenterID, newComment)
}

// DeletePostComment is the resolver for the deletePostComment field.
func (r *mutationResolver) DeletePostComment(ctx context.Context, id string, mainLevel int, subLevel int) (bool, error) {
	return database.DBConnect.DeletePostComment(id, mainLevel, subLevel)
}

// ResetPwd is the resolver for the resetPWD field.
func (r *mutationResolver) ResetPwd(ctx context.Context, password string) (bool, error) {
	id, ok := ctx.Value(string("UserID")).(string)
	if !ok {
		return false, errors.New("failed to get user id from ctx")
	}
	member, err := database.DBConnect.ResetPassword(model.NewPwd{ID: id, Password: password})
	if err != nil {
		log.Error.Print(err)
		return false, err
	}
	err = mail.SendMail(
		member.Email,
		"密碼已被修改 (Your Password Has Been Reset)",
		mail.ResetPWDSuccess(strings.Split(member.Email, "@")[0]),
	)
	if err != nil {
		log.Error.Print(err)
		return false, err
	}
	log.Info.Println("Sent verify mail to " + member.Email)
	return true, err
}

// SendResetPwd is the resolver for the sendResetPWD field.
func (r *mutationResolver) SendResetPwd(ctx context.Context, email string) (bool, error) {
	url := "http://localhost:3000/"
	email = strings.ToLower(email)
	if !database.DBConnect.CheckEmailExist(email) {
		log.Warning.Print("Email is not existed.")
		return false, fmt.Errorf("emailIsNotExisted")
	}
	member, err := database.DBConnect.FindMemberByEmail(email)
	if err != nil {
		return false, fmt.Errorf("emailIsNotInMember")
	}
	token, err := auth.CreateToken(time.Now(), time.Now(), time.Now().Add(time.Duration(10)*time.Minute), member.ID)
	if err != nil {
		log.Warning.Print(err)
		return false, fmt.Errorf("createTokenFailed")
	}
	err = mail.SendMail(
		email,
		"重設密碼 (Reset Password)",
		mail.ResetPWDContent(strings.Split(email, "@")[0], token, url),
	)
	if err != nil {
		log.Error.Print(err)
		return false, fmt.Errorf("sendMailFailed")
	}
	log.Info.Println("Sent ResetPWD mail to " + email)
	return true, nil
}

// UpdateMemberIsAdmin is the resolver for the updateMemberIsAdmin field.
func (r *mutationResolver) UpdateMemberIsAdmin(ctx context.Context, id string) (bool, error) {
	return database.DBConnect.UpdateMemberIsAdmin(id)
}

// SelfInfo is the resolver for the selfInfo field.
func (r *queryResolver) SelfInfo(ctx context.Context) (*model.Member, error) {
	id, ok := ctx.Value(string("UserID")).(string)
	if !ok {
		return nil, errors.New("failed to get user id from ctx")
	}
	return database.DBConnect.FindMemberById(id)
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
	id, ok := ctx.Value(string("UserID")).(string)
	if !ok {
		return false, errors.New("failed to get user id from ctx")
	}
	return database.DBConnect.MemberIsAdmin(id), nil
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	return database.DBConnect.FindPostById(id)
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	return database.DBConnect.AllPost()
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
