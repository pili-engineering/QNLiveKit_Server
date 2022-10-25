package prome

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/qbox/livekit/utils/logger"
	"net/http"
	"time"
)

func Start(ctx context.Context, config Config) error {
	log := logger.ReqLogger(ctx)
	if err := validateConfig(config); err != nil {
		log.Errorf("invalid prome config")
		return err
	}

	if config.ClientMode == ClientModeExporter {
		return startExporter(ctx, config.ExporterConfig)
	} else {
		return startPusher(ctx, config.PusherConfig)
	}
}

func startExporter(ctx context.Context, config ExporterConfig) error {
	log := logger.ReqLogger(ctx)
	http.Handle("/metrics", promhttp.Handler())

	log.Infof("prome exporter start on %s", config.ListenAddr)
	err := http.ListenAndServe(config.ListenAddr, nil)
	log.Errorf("prome exporter stopped, error %s", err.Error())
	return err
}

func startPusher(ctx context.Context, config PusherConfig) error {
	log := logger.ReqLogger(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Errorf("")
			return ctx.Err()

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
}
