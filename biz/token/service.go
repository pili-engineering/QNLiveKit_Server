// @Author: wangsheng
// @Description:
// @File:  service
// @Version: 1.0.0
// @Date: 2022/5/19 9:28 下午
// Copyright 2021 QINIU. All rights reserved

package token

import (
	"time"

	"github.com/qbox/livekit/common/api"

	"github.com/dgrijalva/jwt-go"
)

type AuthToken struct {
	jwt.StandardClaims
	AppId    string
	UserId   string
	DeviceId string
	Role     string
}

var tokenService ITokenService

type Config struct {
	JwtKey string
}

func InitService(conf Config) {
	tokenService = &TokenService{
		jwtKey: conf.JwtKey,
	}
}

func GetService() ITokenService {
	return tokenService
}

type ITokenService interface {
	GenAuthToken(authToken *AuthToken) (string, error)
	ParseAuthToken(token string) (*AuthToken, error)
}

type TokenService struct {
	jwtKey string
}

func (s *TokenService) GenAuthToken(authToken *AuthToken) (string, error) {
	if authToken.ExpiresAt == 0 {
		authToken.ExpiresAt = time.Now().Unix() + 7*86400
	}

	key := s.jwtKey
	signKey := []byte(key)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authToken)
	ss, err := token.SignedString(signKey)
	return ss, err
}

func (s *TokenService) ParseAuthToken(token string) (*AuthToken, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &AuthToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtKey), nil
	})
	if err != nil {
		if isExpiredTokenError(err) {
			return nil, api.ErrTokenExpired
		} else {
			return nil, api.ErrBadToken
		}
	}

	if claims, ok := jwtToken.Claims.(*AuthToken); ok && jwtToken.Valid {
		return claims, nil
	}

	return nil, api.ErrBadToken
}

func isExpiredTokenError(err error) bool {
	if jwtErr, ok := err.(*jwt.ValidationError); ok {
		return jwtErr.Errors&jwt.ValidationErrorExpired != 0
	}
	return false
}
