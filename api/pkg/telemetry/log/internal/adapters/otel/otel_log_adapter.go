package otel

import (
	"frisboo-bank/openapi-generator-service/pkg/telemetry/log/config"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/log/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/log/types/logtype"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

var _ contracts.Log = (*otelLogAdapter)(nil)

type otelLogAdapter struct {
	name   string
	logger loggercontracts.Logger
}

func NewOtelLogAdapter(
	name string,
	cfg *config.LogOptions,
	logger loggercontracts.Logger,
) (contracts.Log, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	return &otelLogAdapter{
		name:   name,
		logger: logger,
	}, nil
}

func (o *otelLogAdapter) Type() logtype.LogType {
	return logtype.LogTypes.OPEN_TELEMETRY
}
func (o *otelLogAdapter) Name() string                   { return o.name }
func (o *otelLogAdapter) Logger() loggercontracts.Logger { return o.logger }
