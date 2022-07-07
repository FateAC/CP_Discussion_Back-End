package database

import (
	"CP_Discussion/env"
	"CP_Discussion/graph/model"
	"context"
	"fmt"
	"log"
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
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		client: client,
	}
}

func (db *DB) InsertMember(input model.NewMember) *model.Member {
	memberColl := db.client.Database(env.DBName).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := memberColl.InsertOne(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	return &model.Member{
		ID:         res.InsertedID.(primitive.ObjectID).Hex(),
		Email:      input.Email,
		Password:   input.Password,
		IsAdmin:    input.IsAdmin,
		Username:   input.Username,
		Nickname:   input.Nickname,
		AvatarPath: input.AvatarPath,
	}
}

func (db *DB) FindMemberById(id string) *model.Member {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	memberColl := db.client.Database(env.DBName).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res := memberColl.FindOne(ctx, bson.M{"_id": ObjectID})
	member := model.Member{}
	res.Decode(&member)
	return &member
}

func (db *DB) AllMember() []*model.Member {
	memberColl := db.client.Database(env.DBName).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := memberColl.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var members []*model.Member
	for cur.Next(ctx) {
		var member *model.Member
		err := cur.Decode(&member)
		if err != nil {
			log.Fatal(err)
		}
		members = append(members, member)
	}
	return members
}

func (db *DB) LoginCheck(input model.Login) *model.Auth {
	email := input.Email
	password := input.Password
	memberColl := db.client.Database(env.DBName).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res := memberColl.FindOne(ctx, bson.M{"email": email})
	member := model.Member{}
	err := res.Decode(&member)
	auth := model.Auth{
		State: false,
		Token: "",
	}

	if err != nil {
		return &auth
	}

	auth.State = member.Password == password

	if auth.State {
		auth.Token = "Hey I am in, please give me a jwt token in the future"
	}

	return &auth
}
