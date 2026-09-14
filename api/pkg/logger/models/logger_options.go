package models

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	encodingtype "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/encoding_type"
	loglevel "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/log_level"
	loggertype "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/logger_type"
	"frisboo-bank/openapi-generator-service/pkg/validation/validators"

	vendorvalidation "github.com/go-ozzo/ozzo-validation"
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
	if l.Level == loglevel.LogLevels.UNKNOWN {
		l.Level = loglevel.LogLevels.INFOLEVEL
	}
	if l.Encoding == encodingtype.EncodingTypes.UNKNOWN {
		l.Encoding = encodingtype.EncodingTypes.JSON
	}
	if l.CallDepth < 0 {
		l.CallDepth = 0
	}
}

func (l *LoggerOptions) Validate() error {
	return vendorvalidation.ValidateStruct(
		l,
		vendorvalidation.Field(&l.Type, vendorvalidation.Required, validators.ValidEnum()),
		vendorvalidation.Field(&l.CallDepth, vendorvalidation.Min(0)),
		vendorvalidation.Field(&l.Encoding, vendorvalidation.Required, validators.ValidEnum()),
		vendorvalidation.Field(&l.Level, vendorvalidation.Required, validators.ValidEnum()),
	)
}
