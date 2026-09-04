package models

import (
	"fmt"
	"strings"

	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	migrationtype "frisboo-bank/openapi-generator-service/pkg/database/migration/models/enums/migration_type"
)

var _ configContracts.Configurable = (*MigrationOptions)(nil)

type MigrationOptions struct {
	IsEnabled     bool                        `mapstructure:"enabled"`
	Type          migrationtype.MigrationType `mapstructure:"type"`
	Debug         bool                        `mapstructure:"debug"`
	MigrationsDir string                      `mapstructure:"migrationsDir"`

	// dependencies
	Logger   string `mapstructure:"logger"`
	DBClient string `mapstructure:"dbClient"`
}

func (c *MigrationOptions) GetEnabled() bool  { return c.IsEnabled }
func (c *MigrationOptions) GetLogger() string { return c.Logger }

func (c *MigrationOptions) SetDefaults() {}

func (c *MigrationOptions) Validate() error {
	if !c.IsEnabled {
		return nil
	}
	if !c.Type.IsValid() {
		return fmt.Errorf("invalid migration type")
	}
	if strings.TrimSpace(c.MigrationsDir) == "" {
		return fmt.Errorf("migrationDir is required")
	}
	if strings.TrimSpace(c.DBClient) == "" {
		return fmt.Errorf("dbClient is required")
	}

	return nil
}
