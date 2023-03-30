// @Author: wangsheng
// @Description:
// @File:  main
// @Version: 1.0.0
// @Date: 2022/5/18 5:14 下午
// Copyright 2021 QINIU. All rights reserved

package main

import (
	"flag"
	"log"
	_ "net/http/pprof"

	"github.com/qbox/livekit/core/application"
	_ "github.com/qbox/livekit/module/base/auth/module"
	_ "github.com/qbox/livekit/module/base/callback"
	_ "github.com/qbox/livekit/module/base/live/module"
	_ "github.com/qbox/livekit/module/base/stats/module"
	_ "github.com/qbox/livekit/module/base/user/module"
	_ "github.com/qbox/livekit/module/biz/censor/module"
	_ "github.com/qbox/livekit/module/biz/gift/module"
	_ "github.com/qbox/livekit/module/biz/item/module"
	_ "github.com/qbox/livekit/module/biz/mic/module"
	_ "github.com/qbox/livekit/module/biz/relay/module"
	_ "github.com/qbox/livekit/module/extend/prome"
	_ "github.com/qbox/livekit/module/fun/im"
	_ "github.com/qbox/livekit/module/fun/pili"
	_ "github.com/qbox/livekit/module/fun/rtc"
	_ "github.com/qbox/livekit/module/store/cache"
	_ "github.com/qbox/livekit/module/store/mysql"
)

var confPath = flag.String("f", "", "live -f /path/to/config")
var confPath2 = flag.String("e", "", "live -e /path/to/config")

func main() {
	flag.Parse()
	err := application.StartWithConfig(*confPath, *confPath2)
	log.Println("application exit ", err.Error())
}
