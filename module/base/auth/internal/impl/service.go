package impl

import (
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/module/base/auth/internal/token"
)

var instance *ServiceImpl

func GetInstance() *ServiceImpl {
	return instance
}

type ServiceImpl struct {
	auth.Config
	tokenService token.ITokenService
}

func ConfigService(conf auth.Config) {
	if conf.JwtKey == "" {
		conf.JwtKey = "jwtkey"
	}
	instance = &ServiceImpl{
		Config:       conf,
		tokenService: token.NewService(conf.JwtKey),
	}
}

func (s *ServiceImpl) RegisterAuthMiddleware() {
	s.RegisterClientAuth()
	s.RegisterServerAuth()
	s.RegisterAdminAuth()
}

func (s *ServiceImpl) GenAuthToken(authToken *token.AuthToken) (string, error) {
	return s.tokenService.GenAuthToken(authToken)
}
