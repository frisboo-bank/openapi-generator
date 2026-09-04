package contracts

import (
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
)

type ConfigLoader interface {
	Load(env environmentEnum.Environment, cfg any) error
	LoadKey(env environmentEnum.Environment, cfg any, key string) error
	HasKey(env environmentEnum.Environment, key string) (bool, error)
}
