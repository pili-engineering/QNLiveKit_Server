// @Author: wangsheng
// @Description:
// @File:  handlefunc
// @Version: 1.0.0
// @Date: 2022/5/20 2:17 下午
// Copyright 2021 QINIU. All rights reserved

package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/module/base/user"
	"github.com/qbox/livekit/utils/logger"
)

func (s *ServiceImpl) RegisterClientAuth() {
	httpq.SetClientAuth(func(ctx *gin.Context) {
		log := logger.ReqLogger(ctx)

		authInfo := ctx.GetHeader("Authorization")
		if len(authInfo) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rest.ErrUnauthorized.WithRequestId(log.ReqID()))
			return
		}

		authToken, err := s.tokenService.ParseAuthToken(authInfo)
		if err != nil {
			switch err1 := err.(type) {
			case *rest.Error:
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, err1.WithRequestId(log.ReqID()))
			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, rest.ErrInternal.WithRequestId(log.ReqID()))
			}
			return
		}

		userService := user.GetService()
		userEntity, err := userService.FindUser(ctx, authToken.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rest.ErrUnauthorized.WithRequestId(log.ReqID()))
			return
		}

		uInfo := auth.UserInfo{
			AppId:    authToken.AppId,
			UserId:   authToken.UserId,
			DeviceId: authToken.DeviceId,
			ImUserId: userEntity.ImUserid,
		}
		ctx.Set(auth.UserCtxKey, &uInfo)

		ctx.Next()
	})
}
