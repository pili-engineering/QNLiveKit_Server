package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/qbox/livekit/utils/logger"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		sourceIP := ctx.ClientIP()
		method := ctx.Request.Method
		requestLog := logger.GetLoggerFromReq(ctx.Request)
		ctx.Set(logger.LoggerCtxKey, requestLog)

		ctx.Writer.Header().Set(logger.RequestIDHeaderKey, requestLog.ReqID())
		ctx.Set(logger.RequestIDHeaderKey, requestLog.ReqID())

		requestLog.LogEntry.WithFields(log.Fields{
			"time_start":              start,
			"method":                  method,
			"source_ip":               sourceIP,
			"path":                    path,
			logger.RequestIDHeaderKey: requestLog.ReqID(),
		}).Infof("start")

		ctx.Next()

		latency := time.Since(start)
		status := ctx.Writer.Status()
		if len(ctx.Errors) > 0 {
			requestLog.LogEntry = requestLog.LogEntry.WithField("error", ctx.Errors.String())
		}
		requestLog.LogEntry.WithFields(log.Fields{
			"time_latency":            fmt.Sprint(latency),
			"status":                  status,
			"time_finish":             time.Now(),
			logger.RequestIDHeaderKey: requestLog.ReqID(),
		}).Infof("finish")
	}
}
