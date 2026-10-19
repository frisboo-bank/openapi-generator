package types

import (
	"reflect"

	"frisboo-bank/openapi-generator-service/pkg/waiter/types/hookstatus"

	"github.com/go-viper/mapstructure/v2"
)

func WaiterEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeFor[hookstatus.HookStatus]():
			return hookstatus.ParseHookStatus(data)
		}

		return data, nil
	}
}
