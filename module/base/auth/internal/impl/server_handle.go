// @Author: wangsheng
// @Description:
// @File:  qiniumac
// @Version: 1.0.0
// @Date: 2022/5/19 3:40 下午
// Copyright 2021 QINIU. All rights reserved

package impl

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/qiniumac"
)

func (s *ServiceImpl) RegisterServerAuth() {
	if !s.Server.Enable {
		return
	}

	httpq.SetServerAuth(func(ctx *gin.Context) {
		log := logger.ReqLogger(ctx)

		auth := ctx.GetHeader("Authorization")
		if len(auth) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rest.ErrUnauthorized.WithRequestId(log.ReqID()))
			return
		}

		if err := s.checkToken(ctx, auth); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rest.ErrUnauthorized.WithRequestId(log.ReqID()))
			return
		}

		ctx.Next()
	})
}

func (s *ServiceImpl) checkToken(ctx *gin.Context, auth string) *rest.Error {
	log := logger.ReqLogger(ctx)

	ak, signExp, err := s.parseToken(ctx, auth)
	if err != nil {
		return err
	}

	if ak != s.Server.AccessKey {
		log.Errorf("access key not match %s", ak)
		return rest.ErrUnauthorized
	}

	sign, err1 := qiniumac.SignRequest([]byte(s.Server.SecretKey), ctx.Request)
	if err1 != nil {
		log.Errorf("sign request error %s", err.Error())
		return rest.ErrInternal
	}

	signStr := base64.URLEncoding.EncodeToString(sign)
	if signStr != signExp {
		log.Errorf("parseAuth: checksum error")
		return rest.ErrUnauthorized
	}

	return nil
}

// 解析token 参数
// 返回：ak，sign，error
func (s *ServiceImpl) parseToken(ctx *gin.Context, auth string) (string, string, *rest.Error) {
	log := logger.ReqLogger(ctx)
	if !strings.HasPrefix(auth, "Qiniu ") {
		log.Errorf("bad token %s", auth)
		return "", "", rest.ErrUnauthorized
	}

	token := auth[6:]
	pos := strings.Index(token, ":")
	if pos == -1 {
		log.Errorf("bad token %s", auth)
		return "", "", rest.ErrUnauthorized
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
