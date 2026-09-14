package enums

import (
	"reflect"

	cachetype "frisboo-bank/openapi-generator-service/pkg/cache/models/enums/cache_type"

	"github.com/go-viper/mapstructure/v2"
)

func CacheEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeFor[cachetype.CacheType]():
			return cachetype.ParseCacheType(data)
		}

		return data, nil
	}
}
