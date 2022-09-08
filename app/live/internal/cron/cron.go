// @Author: wangsheng
// @Description:
// @File:  cron
// @Version: 1.0.0
// @Date: 2022/6/1 2:32 下午
// Copyright 2021 QINIU. All rights reserved

package cron

import (
	"context"
	"github.com/qbox/livekit/app/live/internal/report"
	"github.com/qbox/livekit/common/apimonitor"
	"time"

	"gopkg.in/robfig/cron.v2"

	"github.com/qbox/livekit/biz/live"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
)

//TODO: 多节点情况下，只跑一个
// Run 运行 cronjob
func Run() {
	c := cron.New()

	liveService := live.GetService()
	c.AddFunc("0/3 * * * * ?", func() {
		now := time.Now()
		nowStr := now.Format(timestamp.TimestampFormatLayout)
		log := logger.New("TimeoutLiveRoom")
		log.WithFields(map[string]interface{}{"start": nowStr})

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		liveService.TimeoutLiveRoom(ctx, now)
	})

	c.AddFunc("0/3 * * * * ?", func() {
		now := time.Now()
		nowStr := now.Format(timestamp.TimestampFormatLayout)

		log := logger.New("TimeoutLiveUser")
		log.WithFields(map[string]interface{}{"start": nowStr})

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		liveService.TimeoutLiveUser(ctx, now)
	})

	c.AddFunc("0 0 2 * * ?", func() {
		now := time.Now()
		nowStr := now.Format(timestamp.TimestampFormatLayout)

		log := logger.New("ReportOnlineMessage")
		log.WithFields(map[string]interface{}{"start": nowStr})

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		report.GetService().ReportOnlineMessage(ctx)
	})

	c.AddFunc("0/5 * * * * ?", func() {
		now := time.Now()
		nowStr := now.Format(timestamp.TimestampFormatLayout)

		log := logger.New("ReportApiMonitor")
		log.WithFields(map[string]interface{}{"start": nowStr})

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		apimonitor.ReportMonitorItems(ctx)
	})

	c.Start()
}
