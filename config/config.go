package config

import (
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
	yaml "gopkg.in/yaml.v2"

	"fileServer/utils"
)

// App 全局配置文件实例
var App Config

// Init 读取全局配置文件
func Init() {
	paths := make([]string, 0, 30)
	paths = append(paths, "./config/default.yaml")
	_, err := os.Stat("/etc/le5leFileServer.yaml")
	if err == nil || os.IsExist(err) {
		paths = append(paths, "/etc/le5leFileServer.yaml")
	} else {
		yamls, err := WalkDir("./configs", ".yaml")
		if err == nil {
			paths = append(paths, yamls...)
		}
	}

	for _, c := range paths {
		data, err := ioutil.ReadFile(c)
		if err == nil {
			yaml.Unmarshal(data, &App)
			log.Debug().Caller().Str("func", "config.Init").Msgf("Read config: config=%v, app=%v", c, App)
		} else {
			log.Warn().Err(err).Msgf("Read config error.")
		}
	}

	getEnvConfig()
	if App.Port == 0 {
		App.Port = 8201
	}
	log.Info().Msgf("App config: %v", App)
}

func getEnvConfig() {
	text := os.Getenv("NAME")
	if text != "" {
		App.Name = text
	}

	text = os.Getenv("VERSION")
	if text != "" {
		App.Version = text
	}

	text = os.Getenv("PORT")
	if text != "" {
		App.Port = uint16(utils.Int(text))
	}

	text = os.Getenv("CPU")
	if text != "" {
		App.CPU = utils.Int(text)
	}

	text = os.Getenv("JWT")
	if text != "" {
		App.Jwt = text
	}

	text = os.Getenv("SECRET")
	if text != "" {
		App.Secret = text
	}

	text = os.Getenv("MONGO_ADDRESS")
	if text != "" {
		App.Mongo.Address = text
	}

	text = os.Getenv("MONGO_DATABASE")
	if text != "" {
		App.Mongo.Database = text
	}

	text = os.Getenv("MONGO_USER")
	if text != "" {
		App.Mongo.User = text
	}

	text = os.Getenv("MONGO_PASSWORD")
	if text != "" {
		App.Mongo.Password = text
	}

	text = os.Getenv("MONGO_MAXCONNECTTIONS")
	if text != "" {
		App.Mongo.MaxConnections = utils.Int(text)
	}

	text = os.Getenv("MONGO_TIMEOUT")
	if text != "" {
		App.Mongo.Timeout = utils.Int(text)
	}

	text = os.Getenv("MONGO_MECHANISM")
	if text != "" {
		App.Mongo.Mechanism = text
	}
	text = os.Getenv("MONGO_AUTHSOURCE")
	if text != "" {
		App.Mongo.AuthSource = text
	}
	text = os.Getenv("MONGO_DEBUG")
	if text == "true" {
		App.Mongo.Debug = true
	}

	text = os.Getenv("LOG_FILENAME")
	if text != "" {
		App.Log.Filename = text
	}

	text = os.Getenv("LOG_MAXSIZE")
	if text != "" {
		App.Log.MaxSize = utils.Int(text)
	}

	text = os.Getenv("LOG_MAXBACKUPS")
	if text != "" {
		App.Log.MaxBackups = utils.Int(text)
	}

	text = os.Getenv("LOG_MAXAGE")
	if text != "" {
		App.Log.MaxAge = utils.Int(text)
	}

}
