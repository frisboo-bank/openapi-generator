package enums

import (
	"reflect"

	metricstype "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/models/enums/metrics_type"

	"github.com/go-viper/mapstructure/v2"
)

func MetricsEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeFor[metricstype.MetricsType]():
			return metricstype.ParseMetricsType(data)
		}

		return data, nil
	}
}
