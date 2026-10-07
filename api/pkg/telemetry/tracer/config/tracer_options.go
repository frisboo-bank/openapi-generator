package config

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/types/tracertype"
)

var _ configContracts.Configurable = (*TracerOptions)(nil)

type TracerOptions struct {
	Type tracertype.TracerType `mapstructure:"type" json:"type"`

	// OpenTelemetry
	Endpoint string `mapstructure:"endpoint" json:"endpoint"`
	Insecure bool   `mapstructure:"insecure" json:"insecure"`

	// dependencies
	Logger string `mapstructure:"logger" json:"logger"`
}

func (o *TracerOptions) GetLogger() string { return o.Logger }

func (o *TracerOptions) SetDefaults() {}

func (o *TracerOptions) Validate() error {
	return nil
}
