package live

import (
	"github.com/qbox/livekit/module/base/live/internal/impl"
	"github.com/qbox/livekit/module/base/live/service"
)

func GetService() service.IService {
	return service.Instance
}

func InitService() {
	service.Instance = &impl.Service{}
}
