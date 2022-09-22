package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	_ "github.com/qbox/livekit/app/module-test/modules/httptest"
	"github.com/qbox/livekit/core/application"
	_ "github.com/qbox/livekit/module/store/cache"
	_ "github.com/qbox/livekit/module/store/mysql"
)

var confPath = flag.String("f", "", "live -f /path/to/config")

func main() {
	flag.Parse()

	err := application.StartWithConfig(*confPath)
	log.Printf("application finished, error %v", err)
}