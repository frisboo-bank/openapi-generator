package models

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
	CancelOnShutdownSignal bool   `mapstructure:"cancelOnShutdownSignal"`
	CleanupTimeoutMs       int    `mapstructure:"cleanupTimeoutMs"`
	Logger                 string `mapstructure:"logger"`
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
