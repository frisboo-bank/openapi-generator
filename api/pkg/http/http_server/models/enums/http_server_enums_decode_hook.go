package enums

import (
	"reflect"

	httpservertype "frisboo-bank/openapi-generator-service/pkg/http/http_server/models/enums/http_server_type"

	"github.com/go-viper/mapstructure/v2"
)

func HTTPServerEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeFor[httpservertype.HttpServerType]():
			return httpservertype.ParseHttpServerType(data)
		}

		return data, nil
	}
}
