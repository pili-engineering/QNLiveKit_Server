// @Author: wangsheng
// @Description:
// @File:  route
// @Version: 1.0.0
// @Date: 2022/5/19 2:56 下午
// Copyright 2021 QINIU. All rights reserved

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/app/live/internal/controller/client"
	"github.com/qbox/livekit/app/live/internal/controller/server"
	"github.com/qbox/livekit/utils/logger"
)

func Engine() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery(), logger.LoggerHandleFunc())

	server.RegisterRoute(engine)
	client.RegisterRoute(engine)

	return engine
}
