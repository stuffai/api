package config

import (
	"fmt"
	"net/url"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type obj struct {
	Env      string
	MongoURI string `envconfig:"MONGO_URI"`
}

func (o obj) MongoHost() string {
	u, _ := url.Parse(o.MongoURI)
	return u.Host
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
	if cfg.Env == "" {
		panic("config.init: env.STUFFAI_API_ENV required")
	}
	projectID = fmt.Sprintf("stuffai-%s", cfg.Env)

	log.WithField("mongo", cfg.MongoHost()).Info("config.init: loaded")
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
