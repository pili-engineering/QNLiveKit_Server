package cron

import (
	"context"
	"time"

	"github.com/qbox/livekit/core/module/cron"
	"github.com/qbox/livekit/module/base/live/internal/impl"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
)

//liveService := live.GetService()
//
//	// 定时老化直播间，单节点执行
//	c.AddFunc("0/3 * * * * ?", func() {
//		if !isSingleTaskNode() {
//			return
//		}
//
//		now := time.Now()
//		nowStr := now.Format(timestamp.TimestampFormatLayout)
//		log := logger.New("TimeoutLiveRoom")
//		log.WithFields(map[string]interface{}{"start": nowStr})
//
//		ctx := context.Background()
//		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)
//
//		liveService.TimeoutLiveRoom(ctx, now)
//	})

func RegisterCrons() {
	// 定时老化直播间，单节点执行
	cron.AddSingleTaskFunc("0/3 * * * * ?", func() {
		now := time.Now()
		nowStr := now.Format(timestamp.TimestampFormatLayout)
		log := logger.New("TimeoutLiveRoom")
		log.WithFields(map[string]interface{}{"start": nowStr})

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		impl.GetInstance().TimeoutLiveRoom(ctx, now)
	})

	// 定时老化直播间用户，单节点执行
	cron.AddSingleTaskFunc("0/3 * * * * ?", func() {
		now := time.Now()
		nowStr := now.Format(timestamp.TimestampFormatLayout)

		log := logger.New("TimeoutLiveUser")
		log.WithFields(map[string]interface{}{"start": nowStr})

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		impl.GetInstance().TimeoutLiveUser(ctx, now)
	})

	// 每秒统计缓存中的直播间点赞，写入DB
	// 因为存在补数据，这里不用每秒任务
	cron.AddSingleTaskFunc("@every second", func() {
		log := logger.New("FlushCacheLikes")

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		impl.GetInstance().FlushCacheLikes(ctx)
	})
}
