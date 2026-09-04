package contracts

import "time"

type (
	Fields map[string]any

	Logger interface {
		Debug(args ...any)
		Debugf(template string, args ...any)
		Debugw(msg string, fields Fields)
		Info(args ...any)
		Infof(template string, args ...any)
		Infow(msg string, fields Fields)
		Name() string
		Warn(args ...any)
		Warnf(template string, args ...any)
		WarnMsg(msg string, err error)
		Error(args ...any)
		Errorw(msg string, fields Fields)
		Errorf(template string, args ...any)
		Err(msg string, err error)
		Fatal(args ...any)
		Fatalf(template string, args ...any)
		Printf(template string, args ...any)
		WithName(name string)
		GrpcMiddlewareAccessLogger(
			method string,
			time time.Duration,
			metaData map[string][]string,
			err error,
		)
		GrpcClientInterceptorLogger(
			method string,
			req any,
			reply any,
			time time.Duration,
			metaData map[string][]string,
			err error,
		)
	}
)
