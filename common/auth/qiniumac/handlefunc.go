// @Author: wangsheng
// @Description:
// @File:  qiniumac
// @Version: 1.0.0
// @Date: 2022/5/19 3:40 下午
// Copyright 2021 QINIU. All rights reserved

package qiniumac

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
)

type Config struct {
	Enable    bool   `mapstructure:"enable"`     //是否需要鉴权
	AccessKey string `mapstructure:"access_key"` //AK
	SecretKey string `mapstructure:"secret_key"` //SK
}

type MacAuthMiddleware struct {
	enable    bool   `json:"enable"`     //是否需要鉴权
	accessKey string `json:"access_key"` //A
	secretKey string `json:"secret_key"`
}

func NewAuthMiddleware(conf Config) *MacAuthMiddleware {
	return &MacAuthMiddleware{
		enable:    conf.Enable,
		accessKey: conf.AccessKey,
		secretKey: conf.SecretKey,
	}
}

func (m *MacAuthMiddleware) HandleFunc() gin.HandlerFunc {
	if !m.enable {
		return func(ctx *gin.Context) {
			ctx.Next()
		}
	} else {
		return func(ctx *gin.Context) {
			log := logger.ReqLogger(ctx)

			auth := ctx.GetHeader("Authorization")
			if len(auth) == 0 {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorWithRequestId(log.ReqID(), api.ErrBadToken))
				return
			}

			if err := m.checkToken(ctx, auth); err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorWithRequestId(log.ReqID(), err))
				return
			}

			ctx.Next()
		}
	}
}

func (m *MacAuthMiddleware) checkToken(ctx *gin.Context, auth string) error {
	log := logger.ReqLogger(ctx)

	ak, signExp, err := m.parseToken(ctx, auth)
	if err != nil {
		return err
	}

	if ak != m.accessKey {
		log.Errorf("access key not match %s", ak)
		return api.ErrBadToken
	}

	sign, err := SignRequest([]byte(m.secretKey), ctx.Request)
	if err != nil {
		log.Errorf("sign request error %s", err.Error())
		return api.ErrInternal
	}

	signStr := base64.URLEncoding.EncodeToString(sign)
	if signStr != signExp {
		log.Errorf("parseAuth: checksum error")
		return api.ErrBadToken
	}

	return nil
}

// 解析token 参数
// 返回：ak，sign，error
func (m *MacAuthMiddleware) parseToken(ctx *gin.Context, auth string) (string, string, error) {
	log := logger.ReqLogger(ctx)
	if !strings.HasPrefix(auth, "Qiniu ") {
		log.Errorf("bad token %s", auth)
		return "", "", api.ErrBadToken
	}

	token := auth[6:]
	pos := strings.Index(token, ":")
	if pos == -1 {
		log.Errorf("bad token %s", auth)
		return "", "", api.ErrBadToken
	}

	ak := token[:pos]
	sign := token[pos+1:]
	switch len(sign) % 4 {
	case 3:
		sign += "="
	case 2:
		sign += "=="
	}

	return ak, sign, nil
}
