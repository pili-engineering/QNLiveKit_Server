package appinfo

import (
	"github.com/qbox/livekit/core/module/appinfo/internal/impl"
)

func SetImAppId(imAppId string) {
	impl.SetImAppId(imAppId)
}
