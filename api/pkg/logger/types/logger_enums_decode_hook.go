package types

import (
	"fmt"
	"reflect"
	"strings"

	"frisboo-bank/openapi-generator-service/pkg/logger/types/encodingtype"
	"frisboo-bank/openapi-generator-service/pkg/logger/types/loggertype"
	"frisboo-bank/openapi-generator-service/pkg/logger/types/loglevel"

	"github.com/go-viper/mapstructure/v2"
)

func LoggerEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeFor[encodingtype.EncodingType]():
			return encodingtype.ParseEncodingType(data)
		case reflect.TypeFor[loglevel.LogLevel]():
			str, ok := data.(string)
			if !ok {
				return nil, fmt.Errorf("expected string for log level, got %T", data)
			}
			normalized := strings.ToLower(strings.TrimSpace(str))
			normalized = strings.TrimSuffix(normalized, "level")
			return loglevel.ParseLogLevel(normalized + "Level")
		case reflect.TypeFor[loggertype.LoggerType]():
			return loggertype.ParseLoggerType(data)
		}

		return data, nil
	}
}
