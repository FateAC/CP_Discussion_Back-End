package database

import (
    "time"
    "fmt"
    "context"
    "log"
    "CP_Discussion/graph/model"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
    client *mongo.Client
}

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

func (db *DB) InsertMovieById(movie model.NewMovie) *model.Movie {
    movieColl := db.client.Database("graphql-mongodb-api-db").Collection("movie")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    inserg, err := movieColl.InsertOne(ctx, bson.D{{Key: "name", Value: movie.Name}})

    if err != nil { 
        log.Fatal(err)
    }

    insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
    returnMovie := model.Movie{ID: insertedID, Name: movie.Name}

    return &returnMovie
}

func (db *DB) FindMovieById(id string) *model.Movie {
    ObjectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        log.Fatal(err)
    }

    movieColl := db.client.Database("graphql-mongodb-api-db").Collection("movie")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    res := movieColl.FindOne(ctx, bson.M{"_id": ObjectID})

    movie := model.Movie{ID: id}

    res.Decode(&movie)

    return &movie
}

func (db *DB) All() []*model.Movie {
    movieColl := db.client.Database("graphql-mongodb-api-db").Collection("movie")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    cur, err := movieColl.Find(ctx, bson.D{})
    if err != nil {
        log.Fatal(err)
    }

    var movies []*model.Movie
    for cur.Next(ctx) {
        sus, err := cur.Current.Elements()
        fmt.Println(sus)
        if err != nil {
            log.Fatal(err)
        }

        movie := model.Movie{ID: (sus[0].String()), Name: (sus[1].String())}

        movies = append(movies, &movie)
        }

    return movies
}