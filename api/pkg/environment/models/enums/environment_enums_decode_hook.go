package enums

import (
	"reflect"

	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"

	"github.com/go-viper/mapstructure/v2"
)

func EnvironmentEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		switch t {
		case reflect.TypeFor[environmentEnum.Environment]():
			return environmentEnum.ParseEnvironment(data)
		}

		return data, nil
	}
}
