package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/utils/logger"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		l := logger.GetLoggerFromReq(ctx.Request)
		ctx.Set(logger.LoggerCtxKey, l)
	}
}
