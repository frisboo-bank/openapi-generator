package types

import (
	"reflect"

	logtype "frisboo-bank/openapi-generator-service/pkg/telemetry/log/types/logtype"

	"github.com/go-viper/mapstructure/v2"
)

func LogEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeFor[logtype.LogType]():
			return logtype.ParseLogType(data)
		}

		return data, nil
	}
}
