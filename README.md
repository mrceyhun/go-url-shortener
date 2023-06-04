
### go-url-shortener

[![Go Report Card](https://goreportcard.com/badge/github.com/mrceyhun/go-url-shortener)](https://goreportcard.com/report/github.com/mrceyhun/go-url-shortener)
[![GoDoc](https://godoc.org/github.com/mrceyhun/go-url-shortener?status.svg)](https://godoc.org/github.com/mrceyhun/go-url-shortener)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

Simple web service to shorten long urls and keep hash:url key-value maps in MongoDB.

Service database agnostic and can use any DB.

#### how to

Start a mongo docker container

`$ docker run -it -p 27017:27017 mongo`

Run service

`$ go run . -address "localhost:8080"`

Use either curl for testing:

```
curl -X POST localhost:8080/api/v1/short-url -d '{"Url": "https://host-name/very-long-url?testurl=123&x=y"}'

curl -X GET localhost:8080/api/v1/short-url/c647adf52c439e35daf186bc2a516966
```

Or SwaggerUI:

http://localhost:8080/swagger/index.html#/
