package models

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	encodingtype "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/encoding_type"
	loglevel "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/log_level"
	loggertype "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/logger_type"
)

var _ configContracts.Configurable = (*LoggerOptions)(nil)

type LoggerOptions struct {
	Type          loggertype.LoggerType     `mapstructure:"type"`
	CallDepth     int                       `mapstructure:"callDepth"`
	CallerEnabled bool                      `mapstructure:"callerEnabled"`
	Encoding      encodingtype.EncodingType `mapstructure:"encoding"`
	Level         loglevel.LogLevel         `mapstructure:"level"`
	Prefix        string                    `mapstructure:"prefix"`
	TracerEnabled bool                      `mapstructure:"tracerEnabled"`
}

func (l *LoggerOptions) GetEnabled() bool {
	return true
}

func (l *LoggerOptions) GetLogger() string {
	return ""
}

func (l *LoggerOptions) SetDefaults() {
}

func (l *LoggerOptions) Validate() error {
	return nil
}
