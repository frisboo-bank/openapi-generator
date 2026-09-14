package enums

import (
	"reflect"

	migrationTypeEnum "frisboo-bank/openapi-generator-service/pkg/database/migration/models/enums/migration_type"

	"github.com/go-viper/mapstructure/v2"
)

func MigrationEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		switch t {
		case reflect.TypeFor[migrationTypeEnum.MigrationType]():
			return migrationTypeEnum.ParseMigrationType(data)
		}

		return data, nil
	}
}
