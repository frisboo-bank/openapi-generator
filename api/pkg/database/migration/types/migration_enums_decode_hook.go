package types

import (
	"reflect"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/types/migrationtype"

	"github.com/go-viper/mapstructure/v2"
)

func MigrationEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeFor[migrationtype.MigrationType]():
			return migrationtype.ParseMigrationType(data)
		}

		return data, nil
	}
}
