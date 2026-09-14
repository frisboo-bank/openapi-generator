package enums

import (
	"reflect"

	rpcservertype "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/models/enums/rpc_server_type"

	"github.com/go-viper/mapstructure/v2"
)

func RPCServerEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeFor[rpcservertype.RpcServerType]():
			return rpcservertype.ParseRpcServerType(data)
		}

		return data, nil
	}
}
