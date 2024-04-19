package config

import (
	"net/url"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type obj struct {
	Env       string
	MongoURI  string `envconfig:"MONGO_URI"`
	JWTKey    string `envconfig:"JWT_KEY"`
	ProjectID string `envconfig:"PROJECT_ID"`

	BucketName      string `envconfig:"BUCKET_NAME"`
	TopicIDGenerate string `envconfig:"TOPIC_GENERATE"`
	TopicIDNotify   string `envconfig:"TOPIC_NOTIFY"`
}

func (o obj) MongoHost() string {
	u, _ := url.Parse(o.MongoURI)
	return u.Host
}

var (
	cfg obj
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

	log.WithField("mongo", cfg.MongoHost()).Info("config.init: loaded")
}

func Env() string {
	return cfg.Env
}

func IsLocalEnv() bool {
	return cfg.Env == "local"
}

func BucketName() string {
	return cfg.BucketName
}

func ProjectID() string {
	return cfg.ProjectID
}

func PubSubTopicIDGenerate() string {
	return cfg.TopicIDGenerate
}

func PubSubTopicIDNotify() string {
	return cfg.TopicIDNotify
}

func MongoURI() string {
	return cfg.MongoURI
}

func JWTKey() []byte {
	return []byte(cfg.JWTKey)
}
