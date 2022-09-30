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

	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/admin"
	"github.com/qbox/livekit/module/base/auth/internal/impl"
	token2 "github.com/qbox/livekit/module/base/auth/internal/token"
	"github.com/qbox/livekit/module/base/user"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterRoutes() {
	//authGroup := group.Group("/auth")
	//authGroup.GET("/token", AuthController.GetAuthToken)
	//authGroup.GET("/admin/token", AuthController.GetAdminAuthToken)

	httpq.ServerHandle(http.MethodGet, "/auth/token", AuthController.GetAuthToken)
	httpq.ServerHandle(http.MethodGet, "/auth/admin/token", AuthController.GetAuthToken)
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

type GetAuthTokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

// GetAdminAuthToken 获取管理员的token
func (*authController) GetAdminAuthToken(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)

	req := &GetAuthTokenRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if !req.IsValid() {
		log.Errorf("invalid request %v", req)
		return nil, rest.ErrBadRequest
	}

	adminService := admin.GetManagerService()
	_, err := adminService.FindOrCreateAdminUser(ctx, req.UserId)
	if err != nil {
		log.Errorf("get user userId:%s, error:%v", req.UserId, err)
		return nil, rest.ErrInternal
	}

	authToken := token2.AuthToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: req.ExpiresAt,
		},
		AppId:    req.AppId,
		UserId:   req.UserId,
		DeviceId: req.DeviceId,
		Role:     "admin",
	}

	if token, err := impl.GetInstance().GenAuthToken(&authToken); err != nil {
		log.Errorf("gen token error %v", err)
		return nil, rest.ErrInternal
	} else {
		return &GetAuthTokenResult{
			AccessToken: token,
			ExpiresAt:   authToken.ExpiresAt}, nil
	}
}

// GetAuthToken 获取普通用户的token
func (*authController) GetAuthToken(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)

	req := &GetAuthTokenRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if !req.IsValid() {
		log.Errorf("invalid request %v", req)
		return nil, rest.ErrBadRequest
	}

	userService := user.GetService()
	_, err := userService.FindOrCreateUser(ctx, req.UserId)
	if err != nil {
		log.Errorf("get user userId:%s, error:%v", req.UserId, err)
		return nil, err
	}

	authToken := token2.AuthToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: req.ExpiresAt,
		},
		AppId:    req.AppId,
		UserId:   req.UserId,
		DeviceId: req.DeviceId,
	}

	if token, err := impl.GetInstance().GenAuthToken(&authToken); err != nil {
		log.Errorf("gen token error %v", err)
		return nil, rest.ErrInternal
	} else {
		return &GetAuthTokenResult{
			AccessToken: token,
			ExpiresAt:   authToken.ExpiresAt}, nil
	}
}
