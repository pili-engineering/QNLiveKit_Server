// @Author: wangsheng
// @Description:
// @File:  route
// @Version: 1.0.0
// @Date: 2022/5/19 2:59 下午
// Copyright 2021 QINIU. All rights reserved

package client

import (
	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/module/base/auth/internal/middleware"

	"github.com/qbox/livekit/app/live/internal/config"
	"github.com/qbox/livekit/core/module/appinfo/internal/controller"
)

func RegisterRoute(engine *gin.Engine) {
	clientGroup := engine.Group("/client", middleware.AuthHandleFunc(config.AppConfig.JwtKey))
	controller.RegisterAppRoutes(clientGroup)
	RegisterMicRoutes(clientGroup)
	RegisterStatsRoutes(clientGroup)
}
