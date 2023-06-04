package mongo

// Copyright (c) 2023 - Ceyhun Uzunoglu <ceyhunuzngl AT gmail dot com>

import (
	"context"
	"github.com/mrceyhun/go-url-shortener/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// client MongoDB client
var client *mongo.Client

// dbConnector MongoDB client compatible with model.DbConnector interface
type dbConnector struct {
	collection *mongo.Collection
}

// Insert model.ShortUrl to the MongoDB
func (mc dbConnector) Insert(ctx *context.Context, data *model.ShortUrl) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// if hash is not in db, create new one
	opts := options.Update().SetUpsert(true)
	_, err := mc.collection.UpdateOne(
		*ctx,
		bson.D{{"Hash", data.Hash}},
		bson.D{{"$set", *data}},
		opts)
	return err
}

// FindOne no sort, skip, limit, just match and return model.ShortUrl from MongoDB
func (mc dbConnector) FindOne(ctx *context.Context, hashId string) (model.ShortUrl, error) {
	var result model.ShortUrl
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := mc.collection.FindOne(*ctx, bson.M{"Hash": hashId}, &options.FindOneOptions{}).Decode(&result)
	if err != nil {
		return model.ShortUrl{}, err
	}
	return result, nil
}

// GetMongoDbConnector initialize and return model.DbConnector MongoDB instance
func GetMongoDbConnector(db string, collection string) model.DbConnector {
	return dbConnector{
		collection: client.Database(db).Collection(collection),
	}
}

// ---------------- Connect to MongoDB ----------

// ConnectDb connects to MongoDB
func ConnectDb(mongoUri string, timeout time.Duration) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUri).SetConnectTimeout(timeout)); err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}
}
