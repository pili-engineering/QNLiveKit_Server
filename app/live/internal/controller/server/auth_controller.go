// @Author: wangsheng
// @Description:
// @File:  auth_controller
// @Version: 1.0.0
// @Date: 2022/5/19 2:48 下午
// Copyright 2021 QINIU. All rights reserved

package server

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/biz/token"
	"github.com/qbox/livekit/biz/user"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterAuthRoutes(group *gin.RouterGroup) {
	authGroup := group.Group("/auth")
	authGroup.GET("/token", AuthController.GetAuthToken)
}

var AuthController = &authController{}

type authController struct {
}

type GetAuthTokenRequest struct {
	AppId     string `json:"app_id" form:"app_id"`
	UserId    string `json:"user_id" form:"user_id"`
	DeviceId  string `json:"device_id" form:"device_id"`
	ExpiresAt int64  `json:"expires_at" form:"expires_at"`
}

func (r *GetAuthTokenRequest) IsValid() bool {
	return len(r.AppId) > 0 && len(r.UserId) > 0
}

type GetAuthTokenResponse struct {
	api.Response
	Data struct {
		AccessToken string `json:"access_token"`
		ExpiresAt   int64  `json:"expires_at"`
	} `json:"data"`
}

func (*authController) GetAuthToken(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

	req := &GetAuthTokenRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	if !req.IsValid() {
		log.Errorf("invalid request %v", req)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	userService := user.GetService()
	_, err := userService.FindOrCreateUser(ctx, req.UserId)
	if err != nil {
		log.Errorf("get user userId:%s, error:%v", req.UserId, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), api.ErrInternal))
		return
	}

	authToken := token.AuthToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: req.ExpiresAt,
		},
		AppId:    req.AppId,
		UserId:   req.UserId,
		DeviceId: req.DeviceId,
	}

	tokenService := token.GetService()
	if token, err := tokenService.GenAuthToken(&authToken); err != nil {
		log.Errorf("")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	} else {
		resp := &GetAuthTokenResponse{
			Response: api.Response{
				RequestId: log.ReqID(),
				Code:      0,
				Message:   "success",
			},
		}
		resp.Data.AccessToken = token
		resp.Data.ExpiresAt = authToken.ExpiresAt
		ctx.JSON(http.StatusOK, resp)
	}
}
