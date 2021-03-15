package mongo

import (
	"time"

	"github.com/rs/zerolog/log"
	mgo "gopkg.in/mgo.v2"

	"fileServer/config"
)

// Session 全局mongo会话
var Session *mgo.Session

// Init 初始化mongo连接
func Init() bool {
	var err error
	Session, err = mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:     []string{config.App.Mongo.Address},
		Username:  config.App.Mongo.User,
		Password:  config.App.Mongo.Password,
		Database:  config.App.Mongo.Database,
		Source:    config.App.Mongo.Database,
		Mechanism: config.App.Mongo.Mechanism,
		PoolLimit: config.App.Mongo.MaxConnections,
		Timeout:   time.Duration(config.App.Mongo.Timeout) * time.Second,
	})
	if err == nil {
		if config.App.Mongo.Debug {
			mgo.SetDebug(true)
		}
		if config.App.Mongo.User != "" {
			err = Session.DB(config.App.Mongo.Database).Login(config.App.Mongo.User, config.App.Mongo.Password)
		}
	}

	if err == nil {
		Session.SetMode(mgo.Monotonic, true)
	} else {
		log.Error().Err(err).Msg("Fail to connect mongo.")
		return false
	}

	return true
}
