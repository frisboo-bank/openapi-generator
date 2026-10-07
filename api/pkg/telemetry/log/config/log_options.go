package config

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	logtype "frisboo-bank/openapi-generator-service/pkg/telemetry/log/types/logtype"
)

var _ configContracts.Configurable = (*LogOptions)(nil)

type LogOptions struct {
	IsEnabled bool `mapstructure:"enabled" json:"enabled"`
	Type logtype.LogType `mapstructure:"type" json:"type"`

	// OpenTelemetry
	Endpoint string `mapstructure:"endpoint" json:"endpoint"`
	Insecure bool `mapstructure:"insecure" json:"insecure"`

	// dependencies
	Logger string `mapstructure:"logger" json:"logger"`
}

func (o *LogOptions) Enable() bool  { return o.IsEnabled }
func (o *LogOptions) GetLogger() string { return o.Logger }

func (o *LogOptions) SetDefaults() {}

func (o *LogOptions) Validate() error {
	return nil
}
