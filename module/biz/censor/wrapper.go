package censor

import (
	"github.com/qbox/livekit/module/biz/censor/internal/impl"
	"github.com/qbox/livekit/module/biz/censor/service"
)

func GetService() service.Service {
	return impl.GetInstance()
}
