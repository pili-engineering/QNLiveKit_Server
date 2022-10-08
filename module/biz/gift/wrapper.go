package gift

import (
	"github.com/qbox/livekit/module/biz/gift/internal/impl"
	"github.com/qbox/livekit/module/biz/gift/service"
)

func GetService() service.Service {
	return impl.GetInstance()
}
