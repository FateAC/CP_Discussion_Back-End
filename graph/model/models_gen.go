// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Auth struct {
	Token string `json:"token" bson:"token"`
	State bool   `json:"state" bson:"state"`
}

type Login struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type Member struct {
	ID         string `json:"_id" bson:"_id"`
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	IsAdmin    bool   `json:"isAdmin" bson:"isAdmin"`
	Username   string `json:"username" bson:"username"`
	Nickname   string `json:"nickname" bson:"nickname"`
	AvatarPath string `json:"avatarPath" bson:"avatarPath"`
}

type NewMember struct {
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	IsAdmin    bool   `json:"isAdmin" bson:"isAdmin"`
	Username   string `json:"username" bson:"username"`
	Nickname   string `json:"nickname" bson:"nickname"`
	AvatarPath string `json:"avatarPath" bson:"avatarPath"`
}
