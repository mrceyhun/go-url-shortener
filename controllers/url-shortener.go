package controllers

// Copyright (c) 2023 - Ceyhun Uzunoglu <ceyhunuzngl AT gmail dot com>

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mrceyhun/go-url-shortener/models"
	"net/http"
	"time"
)

// DbClient DB client interface
var DbClient models.DbConnector

// Timeout mongo and gin-gonic context timout
var Timeout time.Duration

// GetUrl request handler to get URL string from given hash of it
func GetUrl(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()
	result, err := DbClient.FindOne(&ctx, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Short url creating error: %s", err.Error()))
	}
	c.JSON(http.StatusOK, result.Url)
}

// CreateShortUrl creates the md5 hash of given URL string and stores it in DB
func CreateShortUrl(c *gin.Context) {
	var req models.ShortUrlReq
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		tempReqBody, _ := c.Get(gin.BodyBytesKey)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Request body: %s", string(tempReqBody.([]byte))))
	}
	hash := utilGetHash(req.Url)
	err = DbClient.Insert(&ctx, &models.ShortUrl{
		Url:  req.Url,
		Hash: hash,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Short url creation failed:%v", err.Error()))
	}
	c.JSON(http.StatusOK, hash)
}

// utilGetHash creates md5 hash of given string
func utilGetHash(u string) string {
	md5Instance := md5.New()
	md5Instance.Write([]byte(u))
	md5Hash := hex.EncodeToString(md5Instance.Sum(nil))
	return md5Hash
}
