package config

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	metricstype "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/types/metricstype"
)

var _ configContracts.Configurable = (*MetricsOptions)(nil)

type MetricsOptions struct {
	Type metricstype.MetricsType `mapstructure:"type" json:"type"`

	// OpenTelemetry
	Endpoint string `mapstructure:"endpoint" json:"endpoint"`
	Insecure bool `mapstructure:"insecure" json:"insecure"`

	// dependencies
	Logger string `mapstructure:"logger" json:"logger"`
}

func (o *MetricsOptions) SetDefaults() {}

func (o *MetricsOptions) Validate() error {
	return nil
}

func (o *MetricsOptions) GetLogger() string { return o.Logger }
