package database

import (
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
		env.DBUsername,
		env.DBPassword,
		env.DBUrl,
		env.DBPort,
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
	memberColl := db.client.Database(env.DBName).Collection("member")
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
	}, nil
}

func (db *DB) FindMemberById(id string) (*model.Member, error) {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	memberColl := db.client.Database(env.DBName).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res := memberColl.FindOne(ctx, bson.M{"_id": ObjectID})
	if err = res.Err(); err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	member := model.Member{}
	res.Decode(&member)
	return &member, nil
}

func (db *DB) AllMember() ([]*model.Member, error) {
	memberColl := db.client.Database(env.DBName).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := memberColl.Find(ctx, bson.D{})
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	var members []*model.Member
	for cur.Next(ctx) {
		member := model.Member{}
		err := cur.Decode(&member)
		if err != nil {
			log.Warning.Print(err)
			return nil, err
		}
		members = append(members, &member)
	}
	return members, nil
}

func (db *DB) LoginCheck(input model.Login) *model.Auth {
	email := input.Email
	password := input.Password
	memberColl := db.client.Database(env.DBName).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	auth := model.Auth{
		State: false,
		Token: "",
	}

	res := memberColl.FindOne(ctx, bson.M{"email": email})
	if err := res.Err(); err != nil {
		log.Warning.Print(err)
		return &auth
	}

	member := model.Member{}
	err := res.Decode(&member)
	if err != nil {
		log.Warning.Print(err)
		return &auth
	}

	auth.State = member.Password == password
	if auth.State {
		auth.Token = "Hey I am in, please give me a jwt token in the future"
	}

	return &auth
}
