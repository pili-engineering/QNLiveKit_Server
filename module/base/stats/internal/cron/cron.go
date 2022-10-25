package cron

import (
	"context"
	"time"

	"github.com/qbox/livekit/core/module/cron"
	"github.com/qbox/livekit/module/base/stats"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
)

func RegisterCrons() {
	// 上报直播间信息，单节点执行
	cron.AddSingleTaskFunc("0 0 2 * * ?", func() {
		now := time.Now()
		nowStr := now.Format(timestamp.TimestampFormatLayout)

		log := logger.New("ReportOnlineMessage")
		log.WithFields(map[string]interface{}{"start": nowStr})

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		stats.GetService().ReportOnlineMessage(ctx)
	})
}
