package main

import (
	"fmt"
	"syscall"

	"github.com/judwhite/go-svc"
	"github.com/sirupsen/logrus"

	"code_game/config"
	"code_game/service"
	"code_game/utils/log"
)

func Init() {
	err := config.Load(".yml")
	if err != nil {
		panic(fmt.Sprintf("load config error:%v", err))
	}
	err = log.InitLogger(config.Log())
	if err != nil {
		panic(fmt.Sprintf("init logger error:%v", err))
	}
}

func main() {
	Init()
	srv := service.NewGameService()
	if err := svc.Run(srv, syscall.SIGINT, syscall.SIGTERM); err != nil {
		logrus.Errorln("service stopped %v", err)
	}
}
