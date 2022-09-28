package impl

import (
	"github.com/qbox/livekit/biz/token"
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/module/base/auth/service"
)

var _ service.Service = &ServiceImpl{}

type ServiceImpl struct {
	auth.Config
	tokenService token.ITokenService
}

func NewService(conf auth.Config) *ServiceImpl {
	if conf.JwtKey == "" {
		conf.JwtKey = "jwtkey"
	}
	return &ServiceImpl{
		Config:       conf,
		tokenService: token.NewService(conf.JwtKey),
	}
}

func (s *ServiceImpl) RegisterAuthMiddleware() {
	s.RegisterClientAuth()
	s.RegisterServerAuth()
	s.RegisterAdminAuth()
}
