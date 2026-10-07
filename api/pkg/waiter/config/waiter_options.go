package config

import (
	"fmt"

	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
)

var _ configContracts.Configurable = (*WaiterOptions)(nil)

const (
	defaultWaitTimeoutMs    = 30000
	defaultCleanupTimeoutMs = 5000
)

type WaiterOptions struct {
	CancelOnShutdownSignal bool `mapstructure:"cancelOnShutdownSignal" json:"cancelOnShutdownSignal"`
	CleanupTimeoutMs int `mapstructure:"cleanupTimeoutMs" json:"cleanupTimeoutMs"`
	Logger string `mapstructure:"logger" json:"logger"`
}

func (c *WaiterOptions) GetEnabled() bool {
	return true
}

func (c *WaiterOptions) GetLogger() string {
	return c.Logger
}

func (c *WaiterOptions) SetDefaults() {
	if c.CleanupTimeoutMs <= 0 {
		c.CleanupTimeoutMs = defaultCleanupTimeoutMs
	}
}

func (c *WaiterOptions) Validate() error {
	if c.CleanupTimeoutMs <= 0 {
		return fmt.Errorf("cleanupTimeoutMs must be positive")
	}
	return nil
}

//go:generate go run github.com/invopop/jsonschema -o schema.json -package config WaiterOptions
