package main

// Copyright (c) 2023 - Ceyhun Uzunoglu <ceyhunuzngl AT gmail dot com>

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/mrceyhun/go-url-shortener/controllers"
	"github.com/mrceyhun/go-url-shortener/docs"
	"github.com/mrceyhun/go-url-shortener/mongo"
	"github.com/mrceyhun/go-url-shortener/server"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)

// Flags
var (
	address         = flag.String("address", "localhost:8080", "listen and serve address")
	mongoConfigFile = flag.String("mongo-conf", "./mongo/config.json", "MongoDB configuration file")
	paramTimeout    = flag.Int("timeout", 10, "timeout seconds")
)

// timeout context timout
var timeout = time.Duration(*paramTimeout) * time.Second

// MainRouter main request router
func MainRouter() http.Handler {
	r := gin.New()
	r.Use(gin.Recovery(), server.MiddlewareReqHandler())
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	shortUrl := v1.Group("/short-url")
	{
		shortUrl.GET("/:id", controllers.GetUrlByHash)
		shortUrl.POST("", controllers.CreateShortUrl)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}

// initializeMongoConnection main DB connection function
func initializeMongoConnection() {
	// get MongoDB config and connect to the DB
	err, mongoConf := mongo.ParseConfig(*mongoConfigFile)
	if err != nil {
		log.Printf("cannot read MongoDB config file %s", err)
	}
	log.Println(mongoConf.String())
	mongo.ConnectDb(mongoConf.Uri, timeout)
	controllers.DbClient = mongo.GetMongoDbConnector(mongoConf.Db, mongoConf.Collection)
}

// main
// @title go-url-shortener API documentation
// @version 1.0.0
// @host localhost:8080
// @BasePath /api/v1
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	// get timout seconds from input parameter and create its duration
	timeout = time.Duration(*paramTimeout) * time.Second
	// set controller timeout
	controllers.Timeout = timeout
	// connect to Mongo
	initializeMongoConnection()
	// create server group
	var g errgroup.Group
	mainServer := &http.Server{
		Addr:         *address,
		Handler:      MainRouter(),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}
	g.Go(func() error {
		return mainServer.ListenAndServe()
	})
	if err := g.Wait(); err != nil {
		log.Printf("[ERROR] server failed %s", err)
	}
}
