package models

import configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"

var _ configContracts.Configurable = (*MediatorOptions)(nil)

type MediatorOptions struct {
	Logger string `mapstructure:"logger"`
}

func (c *MediatorOptions) GetEnabled() bool {
	return true
}

func (c *MediatorOptions) GetLogger() string {
	return c.Logger
}

func (c *MediatorOptions) SetDefaults() {
}

func (c *MediatorOptions) Validate() error {
	return nil
}
