// @Author: wangsheng
// @Description:
// @File:  service_test.go
// @Version: 1.0.0
// @Date: 2022/5/19 10:11 下午
// Copyright 2021 QINIU. All rights reserved

package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenService_GenAuthToken(t1 *testing.T) {
	authToken := AuthToken{
		AppId:  "test-app",
		UserId: "test-user",
	}

	tokenService := NewService("jwtTest")
	token, err := tokenService.GenAuthToken(&authToken)
	assert.Nil(t1, err)
	assert.NotEmpty(t1, token)

	token1, err := tokenService.ParseAuthToken(token)
	assert.Nil(t1, err)
	assert.Equal(t1, authToken.AppId, token1.AppId)
	assert.Equal(t1, authToken.UserId, token1.UserId)
	assert.Equal(t1, authToken.DeviceId, token1.DeviceId)

	authToken.ExpiresAt = time.Now().Unix() - 100
	token, err = tokenService.GenAuthToken(&authToken)
	assert.Nil(t1, err)
	assert.NotEmpty(t1, token)

	token2, err := tokenService.ParseAuthToken(token)
	assert.NotNil(t1, err)
	assert.Nil(t1, token2)
}
