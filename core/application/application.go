package application

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/qbox/livekit/utils/logger"
)

func StartWithConfig(path string) error {
	ctx := context.Background()
	log := logger.ReqLogger(ctx)

	errCh := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
		sig := <-c
		errCh <- fmt.Errorf("signal %s", sig)
	}()

	err := <-errCh
	log.Errorf("application will stop, %v", err)

	return err
}
