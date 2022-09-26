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
	"github.com/qbox/livekit/biz/admin"
	"github.com/qbox/livekit/biz/callback"
	"github.com/qbox/livekit/biz/live"
	"github.com/qbox/livekit/biz/report"
	"github.com/qbox/livekit/biz/token"
	"github.com/qbox/livekit/common/prome"
	"github.com/qbox/livekit/common/trace"
	"github.com/qbox/livekit/core/application"
	"github.com/qbox/livekit/core/module/cron"
	"github.com/qbox/livekit/core/module/uuid"
	"github.com/qbox/livekit/module/fun/im"
	"github.com/qbox/livekit/module/fun/rtc"
	"github.com/qbox/livekit/module/store/cache"
	log "github.com/qbox/livekit/utils/logger"
)

var confPath = flag.String("f", "", "live -f /path/to/config")

func main() {
	flag.Parse()
	application.StartWithConfig(*confPath)

	initAllService()
	uuid.Init(config.AppConfig.NodeID)

	errCh := make(chan error)

	go func() {
		err := prome.Start(context.Background(), config.AppConfig.PromeConfig)
		errCh <- err
	}()

	cron.Run()

	report.GetService().ReportOnlineMessage(context.Background())

	err := <-errCh
	log.StdLog.Fatalf("exit %v", err)
}

func initAllService() {
	token.InitService(token.Config{
		JwtKey: config.AppConfig.JwtKey,
	})
	im.InitService(config.AppConfig.ImConfig)
	rtc.InitService(config.AppConfig.RtcConfig)
	callback.InitService(config.AppConfig.Callback)
	report.InitService()
	trace.InitService(trace.Config{
		IMAppID:    config.AppConfig.ImConfig.AppId,
		RTCAppId:   config.AppConfig.RtcConfig.AppId,
		PiliHub:    config.AppConfig.RtcConfig.Hub,
		AccessKey:  config.AppConfig.RtcConfig.AccessKey,
		SecretKey:  config.AppConfig.RtcConfig.SecretKey,
		ReportHost: config.AppConfig.ReportHost,
	})
	live.InitService(live.Config{
		AccessKey: config.AppConfig.RtcConfig.AccessKey,
		SecretKey: config.AppConfig.RtcConfig.SecretKey,
		PiliHub:   config.AppConfig.RtcConfig.Hub,
	})
	admin.InitJobService(admin.Config{
		AccessKey:      config.AppConfig.RtcConfig.AccessKey,
		SecretKey:      config.AppConfig.RtcConfig.SecretKey,
		CensorCallback: config.AppConfig.CensorCallback,
		CensorBucket:   config.AppConfig.CensorBucket,
	})
	cache.Init(&config.AppConfig.CacheConfig)
}
