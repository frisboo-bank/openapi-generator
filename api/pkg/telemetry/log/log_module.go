package log

import (
	"frisboo-bank/openapi-generator-service/pkg/builder/module"

	"frisboo-bank/openapi-generator-service/pkg/telemetry/log/contracts"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/log/config"
	logenums "frisboo-bank/openapi-generator-service/pkg/telemetry/log/types"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/dig"
)

type LogModuleDependencies struct {
	dig.In
}

var LogModule = module.NewMultiInstancesModule(
	module.MultiInstancesModuleOptions[*config.LogOptions, contracts.Log, LogModuleDependencies]{
		Name:             "telemetry.log",
		ConfigKey:        "telemetry.log",
		ConfigDecodeHook: []mapstructure.DecodeHookFunc{logenums.LogEnumsDecodeHook()},
	},
)
