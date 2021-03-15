package main

import (
	"flag"
	"os"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"fileServer/config"
	"fileServer/db"
	"fileServer/db/mongo"
	"fileServer/router"
)

func main() {
	debug := flag.Bool("debug", false, "Sets log level to debug.")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	// 处理panic
	defer func() {
		if err := recover(); err != nil {
			log.Panic().Msgf("%v", err)
		}
	}()

	// 初始化配置
	config.Init()

	// 设置日志
	if config.App.Log.Filename != "" {
		log.Logger = log.Output(&lumberjack.Logger{
			Filename:   config.App.Log.Filename,
			MaxSize:    config.App.Log.MaxSize, // mb
			MaxBackups: config.App.Log.MaxBackups,
			MaxAge:     config.App.Log.MaxAge, // days
		})
	}

	// 最大cpu使用核心数
	runtime.GOMAXPROCS(config.App.CPU)

	// 数据库连接
	if !db.Init() {
		return
	}
	defer mongo.Session.Close()

	// 监听路由
	router.Listen()
}
