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
	"github.com/qbox/livekit/common/auth/qiniumac"
)

func RegisterRoute(engine *gin.Engine) {
	serverGroup := engine.Group("/server", qiniumac.NewAuthMiddleware(config.AppConfig.MacConfig).HandleFunc())

	RegisterAuthRoutes(serverGroup)
	RegisterUserRoutes(serverGroup)
	RegisterLiveRoutes(serverGroup)

	engine.Any("status", StatusCheckController.CheckStatus)
}
