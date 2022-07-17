package database

import (
	"CP_Discussion/auth"
	"CP_Discussion/env"
	"CP_Discussion/graph/model"
	"CP_Discussion/log"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	client *mongo.Client
}

var DBConnect = Connect(
	fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/",
		env.DBInfo["DBUsername"],
		env.DBInfo["DBPassword"],
		env.DBInfo["DBUrl"],
		env.DBInfo["DBPort"],
	),
)

func Connect(dbUrl string) *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Error.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Error.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error.Fatal(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error.Fatal(err)
	}
	log.Info.Printf("%sconnected database to %s\"%s\"%s", "\u001b[32m", "\u001b[38;5;130m", dbUrl, "\u001b[0m")
	return &DB{
		client: client,
	}
}

func (db *DB) InsertMember(input model.NewMember) (*model.Member, error) {
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := memberColl.InsertOne(ctx, input)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	return &model.Member{
		ID:         res.InsertedID.(primitive.ObjectID).Hex(),
		Email:      input.Email,
		Password:   input.Password,
		IsAdmin:    input.IsAdmin,
		Username:   input.Username,
		Nickname:   input.Nickname,
		AvatarPath: input.AvatarPath,
		Courses:    parseCourses(input.Courses),
	}, nil
}

func (db *DB) FindMemberById(id string) (*model.Member, error) {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	member := model.Member{}
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = memberColl.FindOne(ctx, bson.M{"_id": ObjectID}).Decode(&member)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	return &member, nil
}

func (db *DB) AllMember() ([]*model.Member, error) {
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := memberColl.Find(ctx, bson.D{})
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	var members []*model.Member
	err = cur.All(ctx, &members)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	return members, nil
}

func (db *DB) LoginCheck(input model.Login) *model.Auth {
	email := input.Email
	password := input.Password
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resAuth := model.Auth{
		State: false,
		Token: "",
	}
	member := model.Member{}
	err := memberColl.FindOne(ctx, bson.M{"email": email}).Decode(&member)
	if err != nil {
		log.Warning.Print(err)
		return &resAuth
	}
	if member.Password != password {
		return &resAuth
	}
	token, err := auth.CreatToken(member.ID)
	if err != nil {
		log.Warning.Print(err)
		return &resAuth
	}
	resAuth = model.Auth{
		State: true,
		Token: token,
	}
	return &resAuth
}

func parseCourse(course *model.NewCourse) *model.Course {
	return &model.Course{Name: course.Name}
}

func parseCourses(courses []*model.NewCourse) []*model.Course {
	var res []*model.Course
	for _, course := range courses {
		res = append(res, parseCourse(course))
	}
	return res
}

func (db *DB) AddMemberCourse(id string, input model.NewCourse) (*model.Member, error) {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	course := parseCourse(&input)
	filter := bson.M{"_id": ObjectID}
	update := bson.M{"$addToSet": bson.M{"courses": course}}
	after := options.After
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	member := model.Member{}
	err = memberColl.FindOneAndUpdate(ctx, filter, update, &opts).Decode(&member)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	return &member, nil
}

func (db *DB) RemoveMemberCourse(id string, input model.NewCourse) (*model.Member, error) {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	course := parseCourse(&input)
	filter := bson.M{"_id": ObjectID}
	update := bson.M{"$pull": bson.M{"courses": course}}
	after := options.After
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	member := model.Member{}
	err = memberColl.FindOneAndUpdate(ctx, filter, update, &opts).Decode(&member)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	return &member, nil
}

func (db *DB) MemberIsAdmin(id string) bool {
	member, err := db.FindMemberById(id)
	// no error and member is admin
	return err == nil && member.IsAdmin
}
