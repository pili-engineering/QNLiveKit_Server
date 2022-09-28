// @Author: wangsheng
// @Description:
// @File:  route
// @Version: 1.0.0
// @Date: 2022/5/19 3:00 下午
// Copyright 2021 QINIU. All rights reserved

package server

import (
	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/app/live/internal/config"
	"github.com/qbox/livekit/module/base/auth/internal/middleware"
	"github.com/qbox/livekit/module/base/live/internal/controller/server"
	server2 "github.com/qbox/livekit/module/biz/item/internal/controller/server"
)

func RegisterRoute(engine *gin.Engine) {
	serverGroup := engine.Group("/server", middleware.NewAuthMiddleware(config.AppConfig.MacConfig).HandleFunc())

	RegisterAuthRoutes(serverGroup)
	RegisterUserRoutes(serverGroup)
	server.RegisterLiveRoutes(serverGroup)
	server2.RegisterItemRoutes(serverGroup)

	engine.Any("status", StatusCheckController.CheckStatus)
}
