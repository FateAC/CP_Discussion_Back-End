// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Auth struct {
	Token string `json:"token" bson:"token"`
	State bool   `json:"state" bson:"state"`
}

type Comment struct {
	Commenter string    `json:"commenter" bson:"commenter"`
	Content   string    `json:"content" bson:"content"`
	MainLevel int       `json:"mainLevel" bson:"mainLevel"`
	SubLevel  int       `json:"subLevel" bson:"subLevel"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

type Course struct {
	Name string `json:"name" bson:"name"`
}

type Login struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type Member struct {
	ID         string    `json:"_id" bson:"_id"`
	Email      string    `json:"email" bson:"email"`
	Password   string    `json:"password" bson:"password"`
	IsAdmin    bool      `json:"isAdmin" bson:"isAdmin"`
	Username   string    `json:"username" bson:"username"`
	Nickname   string    `json:"nickname" bson:"nickname"`
	AvatarPath string    `json:"avatarPath" bson:"avatarPath"`
	Courses    []*Course `json:"courses" bson:"courses"`
}

type NewComment struct {
	Commenter string `json:"commenter" bson:"commenter"`
	Content   string `json:"content" bson:"content"`
	MainLevel int    `json:"mainLevel" bson:"mainLevel"`
	SubLevel  int    `json:"subLevel" bson:"subLevel"`
}

type NewCourse struct {
	Name string `json:"name" bson:"name"`
}

type NewMember struct {
	Email      string       `json:"email" bson:"email"`
	Password   string       `json:"password" bson:"password"`
	IsAdmin    bool         `json:"isAdmin" bson:"isAdmin"`
	Username   string       `json:"username" bson:"username"`
	Nickname   string       `json:"nickname" bson:"nickname"`
	AvatarPath string       `json:"avatarPath" bson:"avatarPath"`
	Courses    []*NewCourse `json:"courses" bson:"courses"`
}

type NewPwd struct {
	ID       string `json:"id" bson:"id"`
	Password string `json:"password" bson:"password"`
}

type NewPost struct {
	Title  string   `json:"title" bson:"title"`
	Tags   []string `json:"tags" bson:"tags"`
	MdPath string   `json:"mdPath" bson:"mdPath"`
}

type Post struct {
	ID             string    `json:"_id" bson:"_id"`
	Title          string    `json:"title" bson:"title"`
	Tags           []string  `json:"tags" bson:"tags"`
	MdPath         string    `json:"mdPath" bson:"mdPath"`
	CreateTime     time.Time `json:"createTime" bson:"createTime"`
	LastModifyTime time.Time `json:"lastModifyTime" bson:"lastModifyTime"`
}

type SendResetPassword struct {
	Email string `json:"email" bson:"email"`
}
