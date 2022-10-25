// @Author: wangsheng
// @Description:
// @File:  app_controller
// @Version: 1.0.0
// @Date: 2022/6/17 3:55 下午
// Copyright 2021 QINIU. All rights reserved

package client

import (
	"net/http"

	"github.com/qbox/livekit/core/module/appinfo/internal/impl"
	"github.com/qbox/livekit/core/module/httpq"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes() {
	//appGroup := group.Group("/app")
	//appGroup.GET("/config", AppController.GetConfig)

	httpq.ClientHandle(http.MethodGet, "/app/config", GetConfig)
}

// GetConfig 获取应用全局配置信息
// return *impl.AppInfo
func GetConfig(ctx *gin.Context) (interface{}, error) {
	return impl.GetAppInfo(), nil
}
