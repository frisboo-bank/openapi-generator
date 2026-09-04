package enums

import (
	"reflect"

	sqlclientsslmode "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_ssl_mode"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"

	"github.com/go-viper/mapstructure/v2"
)

func SQLClientEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeFor[sqlclienttype.SqlClientType]():
			return sqlclienttype.ParseSqlClientType(data)
		case reflect.TypeFor[sqlclientsslmode.SqlClientSSLMode]():
			return sqlclientsslmode.ParseSqlClientSSLMode(data)
		}

		return data, nil
	}
}
