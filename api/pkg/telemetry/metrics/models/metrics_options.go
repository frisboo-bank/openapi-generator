package models

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	metricstype "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/models/enums/metrics_type"
)

var _ configContracts.Configurable = (*MetricsOptions)(nil)

type MetricsOptions struct {
	IsEnabled bool                    `mapstructure:"enabled"`
	Type      metricstype.MetricsType `mapstructure:"type"`

	// OpenTelemetry
	Endpoint string `mapstructure:"endpoint"`
	Insecure bool   `mapstructure:"insecure"`

	// dependencies
	Logger string `mapstructure:"logger"`
}

func (o *MetricsOptions) GetEnabled() bool  { return o.IsEnabled }
func (o *MetricsOptions) GetLogger() string { return o.Logger }

func (o *MetricsOptions) SetDefaults() {}

func (o *MetricsOptions) Validate() error {
	return nil
}
