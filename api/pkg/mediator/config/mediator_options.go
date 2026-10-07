package config

import configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"

var _ configContracts.Configurable = (*MediatorOptions)(nil)

type MediatorOptions struct {
	Logger string `mapstructure:"logger" json:"logger"`
}

func (c *MediatorOptions) GetLogger() string {
	return c.Logger
}

func (c *MediatorOptions) SetDefaults() {
}

func (c *MediatorOptions) Validate() error {
	return nil
}

//go:generate go run github.com/invopop/jsonschema -o schema.json -package config MediatorOptions
