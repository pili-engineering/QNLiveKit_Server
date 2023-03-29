package application

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/qbox/livekit/core/config"
	"github.com/qbox/livekit/utils/logger"
)

func StartWithConfig(path string, path2 string) error {
	rand.Seed(time.Now().UnixNano())

	if len(path) == 0 {
		return fmt.Errorf("empty config path")
	}
	if len(path2) == 0 {
		return fmt.Errorf("empty config path2")
	}

	ctx := context.Background()
	log := logger.ReqLogger(ctx)

	errCh := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
		sig := <-c
		errCh <- fmt.Errorf("signal %s", sig)
	}()

	c, err := config.LoadConfig(path, path2)
	if err != nil {
		return fmt.Errorf("load config file error %v", err)
	}

	moduleManager.c = c
	go func() {
		err1 := moduleManager.Start()
		if err1 != nil {
			errCh <- err1
		}
	}()

	err = <-errCh
	log.Errorf("application will stop, %v", err)

	return err
}

var stopOnce sync.Once

func Stop(err error) {
	stopOnce.Do(func() {
		moduleManager.Stop(err)
	})
}
