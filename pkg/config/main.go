package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type obj struct {
	Env      string
	MongoURI string `envconfig:"MONGO_URI"`
}

var (
	cfg       obj
	projectID string
)

func init() {
	var err error
	if err = envconfig.Process("stuffai_api", &cfg); err != nil {
		panic("config.init: " + err.Error())
	}
	projectID = fmt.Sprintf("stuffai-%s", cfg.Env)
}

func Env() string {
	return cfg.Env
}

func ProjectID() string {
	return projectID
}

func PubSubTopicID() string {
	return projectID
}

func MongoURI() string {
	return cfg.MongoURI
}
