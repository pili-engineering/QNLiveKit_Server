// @Author: wangsheng
// @Description:
// @File:  route
// @Version: 1.0.0
// @Date: 2022/5/19 2:56 下午
// Copyright 2021 QINIU. All rights reserved

package controller

import (
	"time"

	"github.com/qbox/livekit/app/live/internal/controller/admin"
	"github.com/qbox/livekit/common/apimonitor"
	"github.com/qbox/livekit/core/module/httpq/middleware"
	"github.com/qbox/livekit/utils/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/app/live/internal/controller/client"
	"github.com/qbox/livekit/app/live/internal/controller/server"
)

func Engine() *gin.Engine {
	engine := gin.New()
	engine.Use(Cors(),
		logger.LoggerHandleFunc(),
		middleware.Middleware(),
		apimonitor.Middleware(),
		gin.Recovery(),
	)

	server.RegisterRoute(engine)
	client.RegisterRoute(engine)
	admin.RegisterRoute(engine)
	return engine
}

func Cors() gin.HandlerFunc {
	c := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowHeaders:    []string{"Content-Type", "Access-Token", "Authorization"},
		MaxAge:          6 * time.Hour,
	}

	return cors.New(c)
}
