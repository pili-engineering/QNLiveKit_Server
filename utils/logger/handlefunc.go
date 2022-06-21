// @Author: wangsheng
// @Description:
// @File:  handlefunc
// @Version: 1.0.0
// @Date: 2022/5/19 6:23 下午
// Copyright 2021 QINIU. All rights reserved

package logger

import "github.com/gin-gonic/gin"

func LoggerHandleFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		l := GetLoggerFromReq(ctx.Request)
		ctx.Set(LoggerCtxKey, l)
	}
}
