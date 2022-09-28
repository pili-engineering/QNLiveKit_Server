package live

import (
	"github.com/qbox/livekit/module/base/live/service"
)

func GetService() service.IService {
	return service.Instance
}
