package main

// Copyright (c) 2023 - Ceyhun Uzunoglu <ceyhunuzngl AT gmail dot com>

import (
    "context"
    "crypto/md5"
    "encoding/hex"
    "flag"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/mrceyhun/go-url-shortener/models"
    "github.com/mrceyhun/go-url-shortener/mongo"
    "golang.org/x/sync/errgroup"
    "log"
    "net/http"
    "time"
)

// Flags
var (
    paramListenAndServeAddress = flag.String("address", "localhost:8080", "listen and serve address")
    paramMongoUri              = flag.String("mongo-uri", "", "MongoDB uri")
    paramMongoDb               = flag.String("mongo-db", "", "MongoDB database")
    paramMongoCollection       = flag.String("mongo-col", "", "MongoDB database")
    paramTimeout               = flag.Int("timeout", 10, "timeout seconds")
)

// timeout mongo and gin-gonic context timout
var timeout = time.Duration(*paramTimeout) * time.Second

// dbClient general key-value DB client
var dbClient models.DatabaseClient

// utilGetHash creates the md5 hash of string
func utilGetHash(u string) string {
    md5Instance := md5.New()
    md5Instance.Write([]byte(u))
    md5Hash := hex.EncodeToString(md5Instance.Sum(nil))
    return md5Hash
}

// handlerUrlFromHash request handler to get URL from given hash
func handlerUrlFromHash(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    err, result := dbClient.FindOne(&ctx, c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, fmt.Sprintf("Short url creating error: %s", err.Error()))
    }
    c.JSON(http.StatusOK, result.Url)
}

// handlerUrlFromHash request handler to create the md5 hash of given URL
func handlerUrlShortener(c *gin.Context) {
    var req models.ShortUrlReq
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    err := c.ShouldBindBodyWith(&req, binding.JSON)
    if err != nil {
        tempReqBody, _ := c.Get(gin.BodyBytesKey)
        c.JSON(http.StatusBadRequest, fmt.Sprintf("Request body: %s", string(tempReqBody.([]byte))))
    }
    hash := utilGetHash(req.Url)
    err = dbClient.Insert(&ctx, &models.ShortUrl{
        Url:  req.Url,
        Hash: hash,
    })
    if err != nil {
        c.JSON(http.StatusBadRequest, fmt.Sprintf("Short url creation failed:%v", err.Error()))
    }
    c.JSON(http.StatusOK, hash)
}

// MainRouter main request router
func MainRouter() http.Handler {
    engine := gin.New()
    engine.Use(gin.Recovery(), middlewareReqHandler())
    e := engine.Group("/api/v1")
    {
        e.GET("/short-url/:id", handlerUrlFromHash)
        e.POST("/short-url", handlerUrlShortener)
    }
    return engine
}

// Serve run service
func Serve() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    flag.Parse()
    timeout = time.Duration(*paramTimeout) * time.Second

    mongo.ConnectDb(*paramMongoUri, timeout)

    dbClient = mongo.GetMongoKeyValueClient(*paramMongoDb, *paramMongoCollection)

    var g errgroup.Group
    mainServer := &http.Server{
        Addr:         *paramListenAndServeAddress,
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

// main
func main() {
    Serve()
}
