package database

import (
	"CP_Discussion/auth"
	"CP_Discussion/env"
	"CP_Discussion/file/fileManager"
	"CP_Discussion/graph/model"
	"CP_Discussion/log"
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
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
	if db.CheckEmailExist(input.Email) {
		log.Warning.Print("Email has already existed.")
		return nil, fmt.Errorf("emailExisted")
	}
	username := strings.Split(input.Email, "@")[0]
	res, err := memberColl.InsertOne(
		ctx,
		struct {
			Email      string          `json:"email" bson:"email"`
			Password   string          `json:"password" bson:"password"`
			IsAdmin    bool            `json:"isAdmin" bson:"isAdmin"`
			Username   string          `json:"username" bson:"username"`
			Nickname   string          `json:"nickname" bson:"nickname"`
			AvatarPath string          `json:"avatarPath" bson:"avatarPath"`
			Courses    []*model.Course `json:"courses" bson:"courses"`
		}{
			Email:      input.Email,
			Password:   input.Password,
			IsAdmin:    input.IsAdmin,
			Username:   username,
			Nickname:   username,
			AvatarPath: "",
			Courses:    parseCourses(input.Courses),
		},
	)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	objectID := res.InsertedID.(primitive.ObjectID)
	return db.FindMemberById(objectID.Hex())
}

func (db *DB) DeleteMember(id string) (*model.Member, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	member := model.Member{}
	err = memberColl.FindOneAndDelete(ctx, bson.M{"_id": objectID}).Decode(&member)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	return &member, nil
}

func (db *DB) CheckEmailExist(input string) bool {
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	member := model.Member{}
	err := memberColl.FindOne(ctx, bson.M{"email": input}).Decode(&member)
	return err == nil
}

func (db *DB) FindMemberById(id string) (*model.Member, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	member := model.Member{}
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = memberColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&member)
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

func comparePWD(DBPWD string, LoginPWD string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(DBPWD), []byte(LoginPWD))
	return err == nil
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
	if comparePWD(member.Password, password) {
		token, err := auth.CreateToken(time.Now(), time.Now(), time.Now().Add(time.Duration(24)*time.Hour), member.ID)
		if err != nil {
			log.Warning.Print(err)
			return &resAuth
		}
		resAuth = model.Auth{
			State: true,
			Token: token,
		}
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
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	course := parseCourse(&input)
	filter := bson.M{"_id": objectID}
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
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	course := parseCourse(&input)
	filter := bson.M{"_id": objectID}
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

func (db *DB) UpdateMemberAvatar(id string, avatar graphql.Upload) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	filename := objectID.Hex() + "." + strings.Split(avatar.ContentType, "/")[1]
	avatarPath := fileManager.BuildAvatarPath(filename)
	avatarUrl := fileManager.BuildAvatarUrl(filename)
	log.Debug.Printf("save avatar: %s\n", avatarPath)
	err = fileManager.SaveFile(avatarPath, avatar.File)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	update := bson.M{"$set": bson.M{"avatarPath": avatarUrl}}
	_, err = memberColl.UpdateByID(ctx, objectID, update)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	return true, nil
}

func (db *DB) DeleteMemberAvatar(id string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	member := model.Member{}
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"avatarPath": ""}}
	err = memberColl.FindOneAndUpdate(ctx, filter, update).Decode(&member)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	Url, err := url.Parse(member.AvatarPath)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	_ = os.Remove(filepath.Join("data", Url.Path))
	return true, nil
}

func (db *DB) UpdateMemberNickname(id string, nickname string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	update := bson.M{"$set": bson.M{"nickname": nickname}}
	_, err = memberColl.UpdateByID(ctx, objectID, update)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	return true, nil
}

func (db *DB) MemberIsAdmin(id string) bool {
	member, err := db.FindMemberById(id)
	// no error and member is admin
	return err == nil && member.IsAdmin
}

func (db *DB) InsertPost(input model.NewPost) (*model.Post, error) {
	postColl := db.client.Database(env.DBInfo["DBName"]).Collection("post")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	doc := struct {
		Poster         string           `json:"poster" bson:"poster"`
		Title          string           `json:"title" bson:"title"`
		Year           int              `json:"year" bson:"year"`
		Semester       int              `json:"semester" bson:"semester"`
		Tags           []string         `json:"tags" bson:"tags"`
		MdPath         string           `json:"mdPath" bson:"mdPath"`
		CreateTime     time.Time        `json:"createTime" bson:"createTime"`
		LastModifyTime time.Time        `json:"lastModifyTime" bson:"lastModifyTime"`
		Comments       []*model.Comment `json:"comments" bson:"comments"`
	}{
		Poster:         input.Poster,
		Title:          input.Title,
		Year:           input.Year,
		Semester:       input.Semester,
		Tags:           input.Tags,
		MdPath:         "",
		CreateTime:     time.Now(),
		LastModifyTime: time.Now(),
		Comments:       []*model.Comment{},
	}
	res, err := postColl.InsertOne(ctx, doc)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	objectID := res.InsertedID.(primitive.ObjectID)
	filename := objectID.Hex() + ".md"
	mdPath := fileManager.BuildPostPath(input.Year, input.Semester, filename)
	mdUrl := fileManager.BuildPostUrl(input.Year, input.Semester, filename)
	err = fileManager.SaveFile(mdPath, input.MdFile.File)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	log.Debug.Printf("save markdown: %s\n", mdPath)
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"mdPath": mdUrl}}
	after := options.After
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	post := model.Post{}
	err = postColl.FindOneAndUpdate(ctx, filter, update, &opts).Decode(&post)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	return &post, nil
}

func (db *DB) DeletePost(id string) (*model.Post, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	postColl := db.client.Database(env.DBInfo["DBName"]).Collection("post")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	post := model.Post{}
	err = postColl.FindOneAndDelete(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	mdUrl, err := url.Parse(post.MdPath)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	_ = os.Remove(filepath.Join("data", mdUrl.Path))
	return &post, nil
}

func (db *DB) FindPostById(id string) (*model.Post, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	postColl := db.client.Database(env.DBInfo["DBName"]).Collection("post")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	post := model.Post{}
	err = postColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	return &post, nil
}

func (db *DB) AllPost() ([]*model.Post, error) {
	postColl := db.client.Database(env.DBInfo["DBName"]).Collection("post")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := postColl.Find(ctx, bson.D{})
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	var posts []*model.Post
	err = cur.All(ctx, &posts)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	return posts, nil
}

func (db *DB) ResetPassword(input model.NewPwd) (*model.Member, error) {
	objectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	password := input.Password
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"password": password}}
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

func (db *DB) FindMemberByEmail(email string) (*model.Member, error) {
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	member := model.Member{}
	err := memberColl.FindOne(ctx, bson.M{"email": email}).Decode(&member)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	return &member, nil
}

func (db *DB) UpdateMemberIsAdmin(id string) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	changeAdmin := !db.MemberIsAdmin(id)
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"isAdmin": changeAdmin}}
	log.Info.Println("change admin!!!")
	_, err = memberColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	log.Info.Println("change admin")
	return true, nil
}
