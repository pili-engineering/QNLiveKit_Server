package user

import (
	"github.com/qbox/livekit/module/base/user/internal/impl"
	"github.com/qbox/livekit/module/base/user/service"
)

func GetService() service.IUserService {
	return impl.GetService()
}
