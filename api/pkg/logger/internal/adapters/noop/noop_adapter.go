package noop

import (
	"time"

	"frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	loggertype "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/logger_type"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

var _ contracts.LoggerAdapter = (*noopAdapter)(nil)

type noopAdapter struct {
	name string
}

func NewNoopAdapter(name string) contracts.LoggerAdapter {
	validation.AssertNotEmpty("name", name)

	return &noopAdapter{
		name: name,
	}
}

func (n *noopAdapter) Debug(...any) {}

func (n *noopAdapter) Debugf(string, ...any) {}

func (n *noopAdapter) Debugw(string, contracts.Fields) {}

func (n *noopAdapter) Info(...any) {}

func (n *noopAdapter) Infof(string, ...any) {}

func (n *noopAdapter) Infow(string, contracts.Fields) {}

func (n *noopAdapter) Warn(...any) {}

func (n *noopAdapter) Warnf(string, ...any) {}

func (n *noopAdapter) WarnMsg(string, error) {}

func (n *noopAdapter) Error(...any) {}

func (n *noopAdapter) Errorf(string, ...any) {}

func (n *noopAdapter) Errorw(string, contracts.Fields) {}

func (n *noopAdapter) Err(string, error) {}

func (n *noopAdapter) Fatal(...any) {}

func (n *noopAdapter) Fatalf(string, ...any) {}

func (n *noopAdapter) Printf(string, ...any) {}

func (n *noopAdapter) WithName(string) {}

func (n *noopAdapter) GrpcMiddlewareAccessLogger(string, time.Duration, map[string][]string, error) {}

func (n *noopAdapter) GrpcClientInterceptorLogger(string, any, any, time.Duration, map[string][]string, error) {
}

func (n *noopAdapter) Name() string                { return n.name }
func (n *noopAdapter) Type() loggertype.LoggerType { return loggertype.LoggerTypes.NOOP }
