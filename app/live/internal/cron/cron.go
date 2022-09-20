// @Author: wangsheng
// @Description:
// @File:  cron
// @Version: 1.0.0
// @Date: 2022/6/1 2:32 下午
// Copyright 2021 QINIU. All rights reserved

package cron

import (
	"context"
	"time"

	"github.com/qbox/livekit/app/live/internal/config"
	"github.com/qbox/livekit/app/live/internal/report"
	"github.com/qbox/livekit/common/apimonitor"

	"gopkg.in/robfig/cron.v2"

	"github.com/qbox/livekit/biz/live"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
)

// Run 运行 cronjob
func Run() {
	c := cron.New()

	liveService := live.GetService()

	// 定时老化直播间，单节点执行
	c.AddFunc("0/3 * * * * ?", func() {
		if !isSingleTaskNode() {
			return
		}

		now := time.Now()
		nowStr := now.Format(timestamp.TimestampFormatLayout)
		log := logger.New("TimeoutLiveRoom")
		log.WithFields(map[string]interface{}{"start": nowStr})

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		liveService.TimeoutLiveRoom(ctx, now)
	})

	// 定时老化直播间用户，单节点执行
	c.AddFunc("0/3 * * * * ?", func() {
		if !isSingleTaskNode() {
			return
		}

		now := time.Now()
		nowStr := now.Format(timestamp.TimestampFormatLayout)

		log := logger.New("TimeoutLiveUser")
		log.WithFields(map[string]interface{}{"start": nowStr})

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		liveService.TimeoutLiveUser(ctx, now)
	})

	// 上报直播间信息，单节点执行
	c.AddFunc("0 0 2 * * ?", func() {
		if !isSingleTaskNode() {
			return
		}

		now := time.Now()
		nowStr := now.Format(timestamp.TimestampFormatLayout)

		log := logger.New("ReportOnlineMessage")
		log.WithFields(map[string]interface{}{"start": nowStr})

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		report.GetService().ReportOnlineMessage(ctx)
	})

	// 上报本节点的API 监控信息，所有节点都要执行
	c.AddFunc("0/5 * * * * ?", func() {
		now := time.Now()
		nowStr := now.Format(timestamp.TimestampFormatLayout)

		log := logger.New("ReportApiMonitor")
		log.WithFields(map[string]interface{}{"start": nowStr})

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		apimonitor.ReportMonitorItems(ctx)
	})

	// 每秒统计缓存中的直播间点赞，写入DB
	// 因为存在补数据，这里不用每秒任务
	if isSingleTaskNode() {
		log := logger.New("FlushCacheLikes")

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		go liveService.FlushCacheLikes(ctx)
	}
	c.Start()
}

func isSingleTaskNode() bool {
	return config.AppConfig.NodeID == config.AppConfig.CronConfig.SingleTaskNode
}
