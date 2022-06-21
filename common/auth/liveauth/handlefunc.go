// @Author: wangsheng
// @Description:
// @File:  handlefunc
// @Version: 1.0.0
// @Date: 2022/5/20 2:17 下午
// Copyright 2021 QINIU. All rights reserved

package liveauth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/biz/token"
	"github.com/qbox/livekit/biz/user"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
)

const UserCtxKey = "UserCtxKey"

type UserInfo struct {
	AppId    string
	UserId   string
	DeviceId string
	ImUserId int64
}

func GetUserInfo(ctx context.Context) *UserInfo {
	log := logger.ReqLogger(ctx)
	i := ctx.Value(UserCtxKey)
	if i == nil {
		return nil
	}

	if t, ok := i.(*UserInfo); ok {
		return t
	} else {
		log.Errorf("%+v not user info", i)
		return nil
	}
}

func AuthHandleFunc(jwtKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.ReqLogger(ctx)

		auth := ctx.GetHeader("Authorization")
		if len(auth) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorWithRequestId(log.ReqID(), api.ErrBadToken))
			return
		}

		tokenService := token.GetService()
		authToken, err := tokenService.ParseAuthToken(auth)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorWithRequestId(log.ReqID(), err))
			return
		}

		userService := user.GetService()
		userEntity, err := userService.FindUser(ctx, authToken.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorWithRequestId(log.ReqID(), err))
			return
		}

		uInfo := UserInfo{
			AppId:    authToken.AppId,
			UserId:   authToken.UserId,
			DeviceId: authToken.DeviceId,
			ImUserId: userEntity.ImUserid,
		}
		ctx.Set(UserCtxKey, &uInfo)

		ctx.Next()
	}
}
