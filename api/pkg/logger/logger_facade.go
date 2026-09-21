package logger

import (
	"time"

	"frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	loggertype "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/logger_type"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

var _ contracts.Logger = (*logger)(nil)

type logger struct {
	adapter contracts.LoggerAdapter
}

func (l *logger) getAdapter() contracts.LoggerAdapter {
	validation.AssertNotNil("adapter", l.adapter)
	return l.adapter
}

func (l *logger) Debug(args ...any) {
	l.getAdapter().Debug(args...)
}

func (l *logger) Debugf(template string, args ...any) {
	l.getAdapter().Debugf(template, args...)
}

func (l *logger) Debugw(msg string, fields contracts.Fields) {
	l.getAdapter().Debugw(msg, fields)
}

func (l *logger) Err(msg string, err error) {
	l.getAdapter().Err(msg, err)
}

func (l *logger) Error(args ...any) {
	l.getAdapter().Error(args...)
}

func (l *logger) Errorf(template string, args ...any) {
	l.getAdapter().Errorf(template, args...)
}

func (l *logger) Errorw(msg string, fields contracts.Fields) {
	l.getAdapter().Errorw(msg, fields)
}

func (l *logger) Fatal(args ...any) {
	l.getAdapter().Fatal(args...)
}

func (l *logger) Fatalf(template string, args ...any) {
	l.getAdapter().Fatalf(template, args...)
}

func (l *logger) GrpcClientInterceptorLogger(
	method string,
	req any,
	reply any,
	time time.Duration,
	metaData map[string][]string,
	err error,
) {
	l.getAdapter().GrpcClientInterceptorLogger(method, req, reply, time, metaData, err)
}

func (l *logger) GrpcMiddlewareAccessLogger(
	method string,
	time time.Duration,
	metaData map[string][]string,
	err error,
) {
	l.getAdapter().GrpcMiddlewareAccessLogger(method, time, metaData, err)
}

func (l *logger) Info(args ...any) {
	l.getAdapter().Info(args...)
}

func (l *logger) Infof(template string, args ...any) {
	l.getAdapter().Infof(template, args...)
}

func (l *logger) Infow(msg string, fields contracts.Fields) {
	l.getAdapter().Infow(msg, fields)
}

func (l *logger) Printf(template string, args ...any) {
	l.getAdapter().Printf(template, args...)
}

func (l *logger) Warn(args ...any) {
	l.getAdapter().Warn(args...)
}

func (l *logger) WarnMsg(msg string, err error) {
	l.getAdapter().WarnMsg(msg, err)
}

func (l *logger) Warnf(template string, args ...any) {
	l.getAdapter().Warnf(template, args...)
}

func (l *logger) WithName(name string) {
	l.getAdapter().WithName(name)
}

func (l *logger) Name() string                { return l.getAdapter().Name() }
func (l *logger) Type() loggertype.LoggerType { return l.getAdapter().Type() }
