package config

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/log/types/logtype"
)

var _ configContracts.Configurable = (*LogOptions)(nil)

type LogOptions struct {
	Type logtype.LogType `mapstructure:"type"`

	// OpenTelemetry
	Endpoint string `mapstructure:"endpoint"`
	Insecure bool   `mapstructure:"insecure"`

	// dependencies
	Logger string `mapstructure:"logger"`
}

func (o *LogOptions) GetEnabled() bool  { return true }
func (o *LogOptions) GetLogger() string { return o.Logger }

func (o *LogOptions) SetDefaults() {}

func (o *LogOptions) Validate() error {
	return nil
}
