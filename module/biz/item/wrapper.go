package item

import (
	"github.com/qbox/livekit/module/biz/item/service"
)

func GetService() service.IItemService {
	return service.Instance
}

func InitService() {

}
