package models

import "context"

// DatabaseClient key-value database client interface
type DatabaseClient interface {
	Insert(ctx *context.Context, data *ShortUrl) error
	FindOne(ctx *context.Context, hashId string) (error, ShortUrl)
}

// ShortUrlReq incoming POST request to create short url
type ShortUrlReq struct {
	Url string `bson:"Url" validate:"required"`
}

// ShortUrl struct is used for key:value couples of unique hash and URL
type ShortUrl struct {
	Url  string `bson:"Url" validate:"required"`
	Hash string `bson:"Hash" validate:"required"`
}
