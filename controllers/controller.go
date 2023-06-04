package controllers

// Copyright (c) 2023 - Ceyhun Uzunoglu <ceyhunuzngl AT gmail dot com>

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mrceyhun/go-url-shortener/model"
	"net/http"
	"time"
)

// DbClient DB client interface
var DbClient model.DbConnector

// Timeout mongo and gin-gonic context timout
var Timeout time.Duration

// GetUrlByHash ... get URL string from given hash of it
// @Summary Get URL string of a hash
// @Description get URL string from given hash of it
// @Tags ShortUrl
// @Param id path string true "Hash String"
// @Success 200 {object} model.ShortUrl
// @Failure 404 {object} object
// @Router /short-url/{id} [get]
func GetUrlByHash(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()
	hash := c.Param("id")
	result, err := DbClient.FindOne(&ctx, hash)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("unable to get url from the given hash: %s", err.Error()))
	}
	c.JSON(http.StatusOK, model.ShortUrl{
		Url:  result.Url,
		Hash: hash,
	})
}

// CreateShortUrl ... creates the md5 hash of given URL string and stores it in DB
// @Summary Creates the md5 hash of given URL string and stores it in DB
// @Description Create new hash:url couple from given url
// @Tags ShortUrl
// @Accept json
// @Param user body model.ShortUrlReq true "Short Url Request"
// @Success 200 {object} model.ShortUrl
// @Failure 400,500 {object} object
// @Router /short-url [post]
func CreateShortUrl(c *gin.Context) {
	var req model.ShortUrlReq
	var shortUrl model.ShortUrl
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		tempReqBody, _ := c.Get(gin.BodyBytesKey)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("unable to bind request body: %s", string(tempReqBody.([]byte))))
	}
	shortUrl = model.ShortUrl{
		Url:  req.Url,
		Hash: utilGetHash(req.Url),
	}
	err = DbClient.Insert(&ctx, &shortUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("unable to create short url :%v", err.Error()))
	}
	c.JSON(http.StatusOK, shortUrl)
}

// utilGetHash creates md5 hash of given string
func utilGetHash(u string) string {
	md5Instance := md5.New()
	md5Instance.Write([]byte(u))
	md5Hash := hex.EncodeToString(md5Instance.Sum(nil))
	return md5Hash
}
