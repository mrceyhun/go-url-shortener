package mongo

// Copyright (c) 2023 - Ceyhun Uzunoglu <ceyhunuzngl AT gmail dot com>

import (
	"github.com/mrceyhun/go-url-shortener/models"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

type Connection struct {
	collection *mongo.Collection
}

func GetMongoKeyValueClient(db string, collection string) models.DatabaseClient {
	return Connection{
		collection: Client.Database(db).Collection(collection),
	}
}

// Insert returns count of query result
func (mc Connection) Insert(ctx *context.Context, data *models.ShortUrl) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	_, err := mc.collection.InsertOne(*ctx, *data)
	return err
}

// FindOne no sort, skip, limit, just match
func (mc Connection) FindOne(ctx *context.Context, hashId string) (error, models.ShortUrl) {
	var result models.ShortUrl
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	opts := options.FindOneOptions{}
	x := mc.collection.FindOne(*ctx, bson.M{"Hash": hashId}, &opts)
	if err := x.Decode(&result); err != nil {
		return err, models.ShortUrl{}
	}
	return nil, result
}

// ---------------- initialize ----------

// ConnectDb connects to MongoDB
func ConnectDb(mongoUri string, timeout time.Duration) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if Client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUri).SetConnectTimeout(timeout)); err != nil {
		log.Fatal(err)
	}
	if err = Client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}
}
