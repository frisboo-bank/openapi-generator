package enums

import (
	"fmt"
	"reflect"
	"strings"

	encodingtype "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/encoding_type"
	loglevel "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/log_level"
	loggertype "frisboo-bank/openapi-generator-service/pkg/logger/models/enums/logger_type"

	"github.com/go-viper/mapstructure/v2"
)

func LoggerEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		switch t {
		case reflect.TypeFor[encodingtype.EncodingType]():
			return encodingtype.ParseEncodingType(data)
		case reflect.TypeFor[loglevel.LogLevel]():
			str, ok := data.(string)
			if !ok {
				return nil, fmt.Errorf("expected string for log level, got %T", data)
			}
			data = strings.TrimSuffix(strings.ToLower(str), "level")
			return loglevel.ParseLogLevel(fmt.Sprintf("%sLevel", data))
		case reflect.TypeFor[loggertype.LoggerType]():
			return loggertype.ParseLoggerType(data)
		}

		return data, nil
	}
}
