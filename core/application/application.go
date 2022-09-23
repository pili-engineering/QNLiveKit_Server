package application

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/qbox/livekit/core/config"
	"github.com/qbox/livekit/utils/logger"
)

func StartWithConfig(path string) error {
	if len(path) == 0 {
		return fmt.Errorf("empty config path")
	}

	ctx := context.Background()
	log := logger.ReqLogger(ctx)

	errCh := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
		sig := <-c
		errCh <- fmt.Errorf("signal %s", sig)
	}()

	c, err := config.LoadConfig(path)
	if err != nil {
		return fmt.Errorf("load config file error %v", err)
	}

	moduleManager.c = c
	go func() {
		moduleManager.Start()
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
