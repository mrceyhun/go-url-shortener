# go-url-shortener

Simple web service to shorten long urls and keep hash:url key-value maps in MongoDB.

Service is abstracted to use any database.

#### 1

docker run -it --network host mongo

#### 2

go run . -address ":8080"

#### 3

curl -X POST localhost:8080/api/v1/short-url -d '{"Url": "https://host-name/very-long-url?testurl=123&x=y"}'

curl -X GET localhost:8080/api/v1/short-url/c647adf52c439e35daf186bc2a516966
