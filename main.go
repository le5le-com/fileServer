package main

import (
	"os"
	"runtime"
	"time"

	"github.com/kardianos/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"fileServer/config"
	"fileServer/db"
	"fileServer/router"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
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

	prg := &program{}
	// 构建服务对象
	s, err := service.New(prg, &service.Config{
		Name:        "fileServer",
		DisplayName: "le5le file service",
		Description: "This is a file service maked by le5le.",
	})
	if err != nil {
		log.Err(err).Msgf("Fail to new service\n")
	}

	if len(os.Args) == 2 {
		// 有命令则执行
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Err(err).Msgf("Fail to exec service cmd\n")
		}
	} else { //否则说明是方法启动了
		err = s.Run()
		if err != nil {
			log.Err(err).Msgf("Fail to start service\n")
		}
	}
}

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) Stop(s service.Service) error {

	return nil
}
func (p *program) run() {
	// 数据库连接
	if !db.Init() {
		return
	}

	// 监听路由
	router.Listen()
}
