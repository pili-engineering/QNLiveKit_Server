// @Author: wangsheng
// @Description:
// @File:  main
// @Version: 1.0.0
// @Date: 2022/5/18 5:14 下午
// Copyright 2021 QINIU. All rights reserved

package main

import (
	"context"
	"flag"
	_ "net/http/pprof"

	"github.com/qbox/livekit/app/live/internal/config"
	"github.com/qbox/livekit/biz/report"
	"github.com/qbox/livekit/biz/token"
	"github.com/qbox/livekit/core/application"
	log "github.com/qbox/livekit/utils/logger"
)

var confPath = flag.String("f", "", "live -f /path/to/config")

func main() {
	flag.Parse()
	application.StartWithConfig(*confPath)

	initAllService()

	errCh := make(chan error)

	report.GetService().ReportOnlineMessage(context.Background())

	err := <-errCh
	log.StdLog.Fatalf("exit %v", err)
}

func initAllService() {
	token.InitService(token.Config{
		JwtKey: config.AppConfig.JwtKey,
	})
	report.InitService()

	//admin.InitJobService(admin.Config{
	//	AccessKey:      config.AppConfig.RtcConfig.AccessKey,
	//	SecretKey:      config.AppConfig.RtcConfig.SecretKey,
	//	CensorCallback: config.AppConfig.CensorCallback,
	//	CensorBucket:   config.AppConfig.CensorBucket,
	//})
}
