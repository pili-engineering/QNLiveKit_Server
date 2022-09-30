package mic

import (
	"github.com/qbox/livekit/module/biz/mic/internal/controller/impl"
	"github.com/qbox/livekit/module/biz/mic/service"
)

func GetService() service.IService {
	return impl.GetInstance()
}
