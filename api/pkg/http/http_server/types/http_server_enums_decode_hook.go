package types

import (
	"reflect"

	"frisboo-bank/openapi-generator-service/pkg/http/http_server/types/httpservertype"

	"github.com/go-viper/mapstructure/v2"
)

func HTTPServerEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeFor[httpservertype.HttpServerType]():
			return httpservertype.ParseHttpServerType(data)
		}

		return data, nil
	}
}
