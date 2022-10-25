package prome

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/qbox/livekit/utils/logger"
)

var instance *Service

type Service struct {
	c *Config
}

func newService(c *Config) *Service {
	return &Service{
		c: c,
	}
}

func (s *Service) Start(ctx context.Context) error {
	if s.c.ClientMode == ClientModeExporter {
		return startExporter(ctx, s.c.ExporterConfig)
	} else {
		return startPusher(ctx, s.c.PusherConfig)
	}
}

func startExporter(ctx context.Context, config ExporterConfig) error {
	go func() {
		log := logger.ReqLogger(ctx)
		http.Handle("/metrics", promhttp.Handler())

		log.Infof("prome exporter start on %s", config.ListenAddr)
		err := http.ListenAndServe(config.ListenAddr, nil)
		log.Errorf("prome exporter stopped, error %s", err.Error())
	}()

	return nil
}

func startPusher(ctx context.Context, config PusherConfig) error {
	go func() {
		log := logger.ReqLogger(ctx)
		for {
			select {
			case <-ctx.Done():
				log.Errorf("prome pusher stopped")

			case <-time.After(time.Duration(config.IntervalS) * time.Second):
				pusher := push.New(config.URL, config.Job)
				err := pusher.Gatherer(prometheus.DefaultGatherer).
					Grouping("instance", config.Instance).
					Add()
				if err != nil {
					log.Errorf("prometheus reporter push failed: %+v", err)
				}
			}
		}
	}()

	return nil
}
