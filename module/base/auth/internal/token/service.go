// @Author: wangsheng
// @Description:
// @File:  impl
// @Version: 1.0.0
// @Date: 2022/5/19 9:28 下午
// Copyright 2021 QINIU. All rights reserved

package token

import (
	"time"

	"github.com/qbox/livekit/core/rest"

	"github.com/dgrijalva/jwt-go"
)

type AuthToken struct {
	jwt.StandardClaims
	AppId    string
	UserId   string
	DeviceId string
	Role     string
}

type ITokenService interface {
	GenAuthToken(authToken *AuthToken) (string, error)
	ParseAuthToken(token string) (*AuthToken, error)
}

type Service struct {
	jwtKey string
}

func NewService(jwtKey string) ITokenService {
	return &Service{
		jwtKey: jwtKey,
	}
}

func (s *Service) GenAuthToken(authToken *AuthToken) (string, error) {
	if authToken.ExpiresAt == 0 {
		authToken.ExpiresAt = time.Now().Unix() + 7*86400
	}

	key := s.jwtKey
	signKey := []byte(key)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authToken)
	ss, err := token.SignedString(signKey)
	return ss, err
}

func (s *Service) ParseAuthToken(token string) (*AuthToken, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &AuthToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtKey), nil
	})
	if err != nil {
		if isExpiredTokenError(err) {
			return nil, rest.ErrTokenExpired
		} else {
			return nil, rest.ErrUnauthorized
		}
	}

	if claims, ok := jwtToken.Claims.(*AuthToken); ok && jwtToken.Valid {
		return claims, nil
	}

	return nil, rest.ErrUnauthorized
}

func isExpiredTokenError(err error) bool {
	if jwtErr, ok := err.(*jwt.ValidationError); ok {
		return jwtErr.Errors&jwt.ValidationErrorExpired != 0
	}
	return false
}
