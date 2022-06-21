// @Author: wangsheng
// @Description:
// @File:  main
// @Version: 1.0.0
// @Date: 2022/5/18 5:14 下午
// Copyright 2021 QINIU. All rights reserved

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qbox/livekit/common/rtc"

	"github.com/qbox/livekit/common/im"

	"github.com/qbox/livekit/biz/token"

	"github.com/qbox/livekit/utils/uuid"

	"github.com/qbox/livekit/app/live/internal/config"
	"github.com/qbox/livekit/app/live/internal/controller"
	"github.com/qbox/livekit/app/live/internal/cron"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/mysql"
	log "github.com/qbox/livekit/utils/logger"
)

var confPath = flag.String("f", "", "live -f /path/to/config")

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	if err := config.LoadConfig(*confPath); err != nil {
		panic(err)
	}
	initAllService()
	mysql.Init(config.AppConfig.Mysqls...)

	errCh := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errCh <- fmt.Errorf("%s", <-c)
	}()

	engine := controller.Engine()
	go func() {
		addr := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
		err := engine.Run(addr)
		errCh <- err
	}()

	modelList := []interface{}{
		&model.LiveEntity{},
		&model.LiveRoomUserEntity{},
		&model.LiveMicEntity{},
		&model.LiveUserEntity{},
		&model.RelaySession{},
	}
	mysql.GetLive("").AutoMigrate(modelList...)
	uuid.Init(config.AppConfig.NodeID)

	cron.Run()

	err := <-errCh
	log.StdLog.Fatalf("exit %v", err)
}

func initAllService() {
	token.InitService(token.Config{
		JwtKey: config.AppConfig.JwtKey,
	})

	im.InitService(config.AppConfig.ImConfig)
	rtc.InitService(config.AppConfig.RtcConfig)
}
