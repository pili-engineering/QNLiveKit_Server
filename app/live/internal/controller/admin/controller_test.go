package admin

import (
	"context"
	"testing"

	"github.com/qbox/livekit/common/im"
	"github.com/qbox/livekit/common/mysql"

	"github.com/qbox/livekit/module/biz/item"
)

func TestCensorController_notifyCensorBlock(t *testing.T) {
	mysql.Init(&mysql.ConfigStructure{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "123456",
		Default:  "live",
		Database: "live",
	},
		&mysql.ConfigStructure{
			Host:     "localhost",
			Port:     3306,
			Username: "root",
			Password: "123456",
			Database: "live",
			Default:  "live",
			ReadOnly: true,
		})
	im.InitService(im.Config{
		AppId:    "cigzypnhoyno",
		Endpoint: "https://s-1-3-s-api.maximtop.cn",
		Token:    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJhcHAiOiJjaWd6eXBuaG95bm8iLCJzdWIiOiIxNTUyIiwiY2x1c3RlciI6MCwicm9sZSI6MiwiaWF0IjoxNjMwNTc0MzEzfQ.kp_SBqnwd9V8w8-C3lB1r64JGYtloMNZoJZAB5kFuRfgQc4b6qEsY6jx_nZIFa_noyNu2R15_WsPDOFuB39idg",
	})
	item.InitService(item.Config{
		PiliHub:   "",
		AccessKey: "",
		SecretKey: "",
	})

	c := &CensorController{}
	c.notifyCensorBlock(context.Background(), "1572032182606635008")
}
