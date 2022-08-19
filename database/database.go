package database

import (
	"CP_Discussion/auth"
	"CP_Discussion/env"
	"CP_Discussion/file/fileManager"
	"CP_Discussion/graph/model"
	"CP_Discussion/log"
	"CP_Discussion/timeformat"
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
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
	formatPostTime(&post)
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
	formatPostTime(&post)
	return &post, nil
}

func (db *DB) UpdatePostFile(id string, file graphql.Upload) (bool, error) {
	post, err := db.FindPostById(id)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	filename := id + ".md"
	mdPath := fileManager.BuildPostPath(post.Year, post.Semester, filename)
	err = fileManager.SaveFile(mdPath, file.File)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	log.Debug.Printf("update markdown: %s\n", mdPath)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	postColl := db.client.Database(env.DBInfo["DBName"]).Collection("post")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"lastModifyTime": time.Now()}}
	err = postColl.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	return true, nil
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
	formatPostTime(&post)
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
	for _, post := range posts {
		formatPostTime(post)
	}
	return posts, nil
}

func formatPostTime(post *model.Post) {
	timeformat.FormatTime(&post.CreateTime)
	timeformat.FormatTime(&post.LastModifyTime)
	for _, comment := range post.Comments {
		timeformat.FormatTime(&comment.Timestamp)
	}
}

func (db *DB) AddPostComment(id string, commmenterID string, newComment model.NewComment) (bool, error) {
	postColl := db.client.Database(env.DBInfo["DBName"]).Collection("post")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	commterObjectID, err := primitive.ObjectIDFromHex(commmenterID)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	comment := struct {
		Commenter primitive.ObjectID `json:"commenter" bson:"commenter"`
		Content   string             `json:"content" bson:"content"`
		MainLevel int                `json:"mainLevel" bson:"mainLevel"`
		SubLevel  int                `json:"subLevel" bson:"subLevel"`
		Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
		Deleted   bool               `json:"deleted" bson:"deleted"`
	}{
		Commenter: commterObjectID,
		Content:   newComment.Content,
		MainLevel: newComment.MainLevel,
		SubLevel:  newComment.SubLevel,
		Timestamp: time.Now(),
		Deleted:   false,
	}
	update := bson.M{
		"$push": bson.M{
			"comments": bson.M{
				"$each": bson.A{
					comment,
				},
				"$sort": bson.D{
					{Key: "mainLevel", Value: 1},
					{Key: "subLevel", Value: 1},
				},
			},
		},
	}
	postColl.UpdateByID(
		ctx,
		objectID,
		update,
	)
	return true, nil
}

func (db *DB) DeletePostComment(id string, commmenterID string, mainLevel int, subLevel int) (bool, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	commenterID, err := primitive.ObjectIDFromHex(commmenterID)
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	postColl := db.client.Database(env.DBInfo["DBName"]).Collection("post")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	filter := bson.M{
		"_id": objectID,
		"comments": bson.M{
			"$elemMatch": bson.M{
				"commenter": commenterID,
				"mainLevel": mainLevel,
				"subLevel":  subLevel,
				"deleted":   false,
			},
		},
	}
	update := bson.M{"$set": bson.M{"comments.$[item].deleted": true}}
	arrayFilter := bson.M{
		"item.commenter": commenterID,
		"item.mainLevel": mainLevel,
		"item.subLevel":  subLevel,
		"item.deleted":   false,
	}
	opts := options.FindOneAndUpdate().SetArrayFilters(
		options.ArrayFilters{
			Filters: []interface{}{
				arrayFilter,
			},
		},
	)
	err = postColl.FindOneAndUpdate(ctx, filter, update, opts).Err()
	if err != nil {
		log.Warning.Print(err)
		return false, err
	}
	return true, nil
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

func (db *DB) GetPostsByTags(year int, semester int, tags []string) ([]*model.Post, error) {
	postColl := db.client.Database(env.DBInfo["DBName"]).Collection("post")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	log.Info.Println(bson.D{{Key: "$all", Value: tags}})
	filter := bson.M{"year": year, "semester": semester, "tags": bson.D{{Key: "$all", Value: tags}}}
	opts := options.Find().SetSort(bson.D{{Key: "createTime", Value: 1}})
	matchPost, err := postColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var posts []*model.Post
	err = matchPost.All(ctx, &posts)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	return posts, nil
}

func (db *DB) AllCourses() ([]*model.Course, error) {
	memberColl := db.client.Database(env.DBInfo["DBName"]).Collection("member")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := memberColl.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			{
				Key: "$match",
				Value: bson.D{
					{
						Key:   "isAdmin",
						Value: true,
					},
				},
			},
		},
		bson.D{
			{
				Key: "$group",
				Value: bson.D{
					{
						Key:   "_id",
						Value: "",
					},
					{
						Key: "courses",
						Value: bson.D{
							{
								Key:   "$push",
								Value: "$courses",
							},
						},
					},
				},
			},
		},
		bson.D{
			{
				Key: "$project",
				Value: bson.D{
					{
						Key:   "_id",
						Value: 0,
					},
					{
						Key: "courses",
						Value: bson.D{
							{
								Key: "$reduce",
								Value: bson.D{
									{
										Key:   "input",
										Value: "$courses",
									},
									{
										Key:   "initialValue",
										Value: bson.A{},
									},
									{
										Key: "in",
										Value: bson.D{
											{
												Key: "$setUnion",
												Value: bson.A{
													"$$value",
													"$$this",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{
				Key: "$unwind",
				Value: bson.D{
					{
						Key:   "path",
						Value: "$courses",
					},
				},
			},
		},
		bson.D{
			{
				Key: "$project",
				Value: bson.D{
					{
						Key:   "name",
						Value: "$courses.name",
					},
				},
			},
		},
	})
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	var courses []*model.Course
	err = cursor.All(ctx, &courses)
	if err != nil {
		log.Warning.Print(err)
		return nil, err
	}
	sort.Slice(
		courses,
		func(i, j int) bool {
			parse := func(course *model.Course) (string, string) {
				parts := strings.Split(course.Name, "_")
				year := parts[0]
				semester := parts[1]
				return year, semester
			}
			years := [2]string{}
			semesters := [2]string{}
			years[0], semesters[0] = parse(courses[i])
			years[1], semesters[1] = parse(courses[j])
			if years[0] != years[1] {
				return years[0] < years[1]
			} else if semesters[0] != semesters[1] {
				return semesters[0] > semesters[1]
			}
			return false
		},
	)
	return courses, nil
}
