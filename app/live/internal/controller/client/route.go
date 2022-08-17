// @Author: wangsheng
// @Description:
// @File:  route
// @Version: 1.0.0
// @Date: 2022/5/19 2:59 下午
// Copyright 2021 QINIU. All rights reserved

package client

import (
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/app/live/internal/config"
	"github.com/qbox/livekit/common/auth/liveauth"
)

func RegisterRoute(engine *gin.Engine) {
	clientGroup := engine.Group("/client", liveauth.AuthHandleFunc(config.AppConfig.JwtKey))
	RegisterAppRoutes(clientGroup)
	RegisterUserRoutes(clientGroup)
	RegisterRelayRoutes(clientGroup)
	RegisterLiveRoutes(clientGroup)
	RegisterMicRoutes(clientGroup)
	RegisterItemRoutes(clientGroup)
	RegisterStatsRoutes(clientGroup)
}
