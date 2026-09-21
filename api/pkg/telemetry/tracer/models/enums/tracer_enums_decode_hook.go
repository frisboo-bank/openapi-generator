package enums

import (
	"reflect"

	tracertype "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/models/enums/tracer_type"

	"github.com/go-viper/mapstructure/v2"
)

func TracerEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeFor[tracertype.TracerType]():
			return tracertype.ParseTracerType(data)
		}

		return data, nil
	}
}
