package config

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/types/metricstype"
)

var _ configContracts.Configurable = (*MetricsOptions)(nil)

type MetricsOptions struct {
	Type metricstype.MetricsType `mapstructure:"type"`

	// OpenTelemetry
	Endpoint string `mapstructure:"endpoint"`
	Insecure bool   `mapstructure:"insecure"`

	// dependencies
	Logger string `mapstructure:"logger"`
}

func (o *MetricsOptions) GetEnabled() bool  { return true }
func (o *MetricsOptions) GetLogger() string { return o.Logger }

func (o *MetricsOptions) SetDefaults() {}

func (o *MetricsOptions) Validate() error {
	return nil
}
