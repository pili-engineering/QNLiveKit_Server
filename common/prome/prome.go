package prome

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qbox/livekit/utils/logger"
	"net/http"
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

	log.Infof("prome start on %s", config.ListenAddr)
	return http.ListenAndServe(config.ListenAddr, nil)
}

func startPusher(ctx context.Context, config PusherConfig) error {
	return nil
}
