package zerolog

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/logger/models"
	encodingtype "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/encoding_type"
	loglevel "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/log_level"

	"github.com/rs/zerolog"
)

var _ contracts.Logger = (*zerologAdapter)(nil)

type zerologAdapter struct {
	name   string
	logger *zerolog.Logger

	callDepth     int
	callerEnabled bool
	encoding      encodingtype.EncodingType
	level         loglevel.LogLevel
	prefix        string
}

// levelMapping converts our internal log level to zerolog level
var levelMapping = map[loglevel.LogLevel]zerolog.Level{
	loglevel.LogLevels.DEBUGLEVEL: zerolog.DebugLevel,
	loglevel.LogLevels.INFOLEVEL:  zerolog.InfoLevel,
	loglevel.LogLevels.WARNLEVEL:  zerolog.WarnLevel,
	loglevel.LogLevels.ERRORLEVEL: zerolog.ErrorLevel,
	loglevel.LogLevels.FATALLEVEL: zerolog.FatalLevel,
	loglevel.LogLevels.PANICLEVEL: zerolog.PanicLevel,
	loglevel.LogLevels.TRACELEVEL: zerolog.TraceLevel,
}

func NewZerologAdapter(
	name string,
	cfg *models.LoggerOptions,
	env environmentEnum.Environment,
) *zerologAdapter {
	adapter := &zerologAdapter{
		name:          name,
		callDepth:     cfg.CallDepth,
		callerEnabled: cfg.CallerEnabled,
		encoding:      cfg.Encoding,
		level:         cfg.Level,
		prefix:        cfg.Prefix,
	}
	adapter.configure()
	return adapter
}

func (z *zerologAdapter) configure() {
	var output io.Writer = os.Stdout
	if z.encoding == encodingtype.EncodingTypes.TEXT {
		output = zerolog.ConsoleWriter{
			Out: output,
		}
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	logger = logger.Level(zerolog.InfoLevel)
	if lvl, ok := levelMapping[z.level]; ok {
		logger = logger.Level(lvl)
	}

	if strings.TrimSpace(z.prefix) != "" {
		logger = logger.With().Str("prefix", z.prefix).Logger()
	}

	if z.callerEnabled {
		callDepth := max(z.callDepth, 0)
		logger = logger.With().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + callDepth).Logger()
	}

	z.logger = &logger
}

func (z *zerologAdapter) Debug(args ...any) {
	z.logger.Debug().Msg(fmt.Sprint(args...))
}

func (z *zerologAdapter) Debugf(template string, args ...any) {
	z.logger.Debug().Msgf(template, args...)
}

func (z *zerologAdapter) Debugw(msg string, fields contracts.Fields) {
	evt := z.logger.Debug()
	for k, v := range fields {
		evt = evt.Interface(k, v)
	}
	evt.Msg(msg)
}

func (z *zerologAdapter) Info(args ...any) {
	z.logger.Info().Msg(fmt.Sprint(args...))
}

func (z *zerologAdapter) Infof(template string, args ...any) {
	z.logger.Info().Msgf(template, args...)
}

func (z *zerologAdapter) Infow(msg string, fields contracts.Fields) {
	evt := z.logger.Info()
	for k, v := range fields {
		evt = evt.Interface(k, v)
	}
	evt.Msg(msg)
}

func (z *zerologAdapter) Warn(args ...any) {
	z.logger.Warn().Msg(fmt.Sprint(args...))
}

func (z *zerologAdapter) Warnf(template string, args ...any) {
	z.logger.Warn().Msgf(template, args...)
}

func (z *zerologAdapter) WarnMsg(msg string, err error) {
	z.logger.Warn().Err(err).Msg(msg)
}

func (z *zerologAdapter) Error(args ...any) {
	z.logger.Error().Msg(fmt.Sprint(args...))
}

func (z *zerologAdapter) Errorf(template string, args ...any) {
	z.logger.Error().Msgf(template, args...)
}

func (z *zerologAdapter) Errorw(msg string, fields contracts.Fields) {
	evt := z.logger.Error()
	for k, v := range fields {
		evt = evt.Interface(k, v)
	}
	evt.Msg(msg)
}

func (z *zerologAdapter) Err(msg string, err error) {
	z.logger.Error().Err(err).Msg(msg)
}

func (z *zerologAdapter) Fatal(args ...any) {
	z.logger.Fatal().Msg(fmt.Sprint(args...))
}

func (z *zerologAdapter) Fatalf(template string, args ...any) {
	z.logger.Fatal().Msgf(template, args...)
}

func (z *zerologAdapter) Printf(template string, args ...any) {
	z.logger.Info().Msgf(template, args...)
}

func (z *zerologAdapter) WithName(name string) {
	*z.logger = z.logger.With().Str("logger", name).Logger()
}

func (z *zerologAdapter) GrpcMiddlewareAccessLogger(
	method string,
	dur time.Duration,
	metaData map[string][]string,
	err error,
) {
	evt := z.logger.Info().
		Str("method", method).
		Dur("duration", dur).
		Interface("metadata", metaData)

	if err != nil {
		evt.Err(err)
	}
	evt.Msg("gRPC request completed")
}

func (z *zerologAdapter) GrpcClientInterceptorLogger(
	method string,
	req any,
	reply any,
	dur time.Duration,
	metaData map[string][]string,
	err error,
) {
	evt := z.logger.Info().
		Str("method", method).
		Interface("request", req).
		Interface("reply", reply).
		Dur("duration", dur).
		Interface("metadata", metaData)

	if err != nil {
		evt.Err(err)
	}
	evt.Msg("gRPC client call completed")
}

func (z *zerologAdapter) Name() string {
	return z.name
}
