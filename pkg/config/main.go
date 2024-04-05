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
	JWTKey   string `envconfig:"JWT_KEY"`
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
	if cfg.JWTKey == "" {
		panic("config.init: env.STUFFAI_API_JWT_KEY required")
	}
	projectID = fmt.Sprintf("stuffai-%s", cfg.Env)

	log.WithField("mongo", cfg.MongoHost()).Info("config.init: loaded")
}

func Env() string {
	return cfg.Env
}

func IsLocalEnv() bool {
	return cfg.Env == "local"
}

func ProjectID() string {
	return projectID
}

func PubSubTopicIDGenerate() string {
	return "generate"
}

func MongoURI() string {
	return cfg.MongoURI
}

func JWTKey() []byte {
	return []byte(cfg.JWTKey)
}
