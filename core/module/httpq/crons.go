package httpq

import (
	"context"
	"time"

	"github.com/qbox/livekit/core/module/cron"
	"github.com/qbox/livekit/core/module/httpq/monitor"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
)

func registerMonitorTask() {
	cron.AddFunc("0/5 * * * * ?", func() {
		now := time.Now()
		nowStr := now.Format(timestamp.TimestampFormatLayout)

		log := logger.New("ReportApiMonitor")
		log.WithFields(map[string]interface{}{"start": nowStr})

		ctx := context.Background()
		ctx = context.WithValue(ctx, logger.LoggerCtxKey, log)

		monitor.ReportMonitorItems(ctx)
	})
}
