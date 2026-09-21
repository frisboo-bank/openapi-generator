package models

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	tracertype "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models/enums/tracer_type"
)

var _ configContracts.Configurable = (*TracerOptions)(nil)

type TracerOptions struct {
	IsEnabled bool                  `mapstructure:"enabled"`
	Type      tracertype.TracerType `mapstructure:"type"`

	// OTLP exporter
	Endpoint string `mapstructure:"endpoint"`
	Insecure bool   `mapstructure:"insecure"`

	// dependencies
	Logger string `mapstructure:"logger"`
}

func (o *TracerOptions) GetEnabled() bool  { return o.IsEnabled }
func (o *TracerOptions) GetLogger() string { return o.Logger }

func (o *TracerOptions) SetDefaults() {}

func (o *TracerOptions) Validate() error {
	return nil
}
