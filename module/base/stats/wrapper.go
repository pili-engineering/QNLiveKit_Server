package stats

import (
	"github.com/qbox/livekit/module/base/stats/service"
)

func GetService() service.IService {
	return service.Instance
}
