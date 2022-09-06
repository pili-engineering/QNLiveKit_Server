package prome

import (
	"fmt"
)

type PusherConfig struct {
	URL       string `mapstructure:"url" validate:"nonzero"`
	Job       string `mapstructure:"job" validate:"nonzero"`
	Instance  string `mapstructure:"instance" validate:"nonzero"`
	IntervalS int    `mapstructure:"interval_s" validate:"nonzero"`
}

type ExporterConfig struct {
	ListenAddr string `mapstructure:"listen_addr"`
}

const (
	ClientModeExporter = "exporter"
	ClientModePusher   = "pusher"
)

type Config struct {
	ClientMode     string         `mapstructure:"client_mode"`
	ExporterConfig ExporterConfig `mapstructure:"exporter_config"`
	PusherConfig   PusherConfig   `mapstructure:"pusher_config"`
}

func validateConfig(config Config) error {
	if config.ClientMode == ClientModeExporter {
		return validateExporterConfig(config.ExporterConfig)
	} else if config.ClientMode == ClientModePusher {
		return validatePusherConfig(config.PusherConfig)
	}

	return fmt.Errorf("invalid prome client_mode %s", config.ClientMode)
}

func validateExporterConfig(config ExporterConfig) error {
	if len(config.ListenAddr) == 0 {
		return fmt.Errorf("export_config with empty listen_addr")
	}

	return nil
}

func validatePusherConfig(config PusherConfig) error {
	if len(config.Job) == 0 {
		return fmt.Errorf("push_config with empty job ")
	}
	if len(config.URL) == 0 {
		return fmt.Errorf("push_config with empty url")
	}
	if len(config.Instance) == 0 {
		return fmt.Errorf("push_config with empty instance ")
	}
	if config.IntervalS == 0 {
		return fmt.Errorf("push_config with empty interval_s")
	}
	return nil
}
