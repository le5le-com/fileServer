package mongo

import (
	"context"
	"fileServer/config"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// Init 初始化mongo连接
func Init() bool {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := "mongodb://"
	if config.App.Mongo.User != "" {
		uri = uri + config.App.Mongo.User + ":" + config.App.Mongo.Password + "@"
	}
	uri = uri + config.App.Mongo.Address + "/" + config.App.Mongo.Database
	if config.App.Mongo.User != "" {
		if config.App.Mongo.AuthSource != "" {
			uri = uri + "?authSource=" + config.App.Mongo.AuthSource + "&authMechanism=" + config.App.Mongo.Mechanism
		} else {
			uri = uri + "?authMechanism=" + config.App.Mongo.Mechanism
		}
	}
	clientOptions := options.Client().ApplyURI(uri).SetMaxPoolSize(uint64(config.App.Mongo.MaxConnections))
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error().Err(err).Msgf("Fail to connect mongo: %v", uri)
		return false
	}
	// 检查连接
	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Error().Err(err).Msgf("Ping mongo: %v", uri)
		return false
	}

	return true
}
