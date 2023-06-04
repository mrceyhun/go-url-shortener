package mongo

// Copyright (c) 2023 - Ceyhun Uzunoglu <ceyhunuzngl AT gmail dot com>

import (
	"encoding/json"
	"log"
	"os"
)

// Config MongoDB credentials struct
type Config struct {
	Uri        string // MongoDB URI: mongodb://127.0.0.1:27017/?directConnection=true
	Db         string // MongoDB database name to use: test
	Collection string // Collection name to use: url-shortener
}

// String marshal the Config struct for pretty string representation
func (c *Config) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		log.Println("unable to marshal Mongo Config", err)
		return ""
	}
	return string(data)
}

// ParseConfig parses given configuration file and initialize Config object
func ParseConfig(configFile string) (error, Config) {
	var config Config
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Println("unable to read config file", configFile, err)
		return err, config
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("unable to parse config file", configFile, err)
		return err, config
	}
	return nil, config
}
