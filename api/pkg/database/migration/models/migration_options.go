package models

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	migrationtype "frisboo-bank/openapi-generator-service/pkg/database/migration/models/enums/migration_type"
	"frisboo-bank/openapi-generator-service/pkg/validation/validators"

	vendorvalidation "github.com/go-ozzo/ozzo-validation"
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
	return vendorvalidation.ValidateStruct(
		c,
		vendorvalidation.Field(&c.Type, vendorvalidation.Required, validators.ValidEnum()),
		vendorvalidation.Field(&c.MigrationsDir, vendorvalidation.Required),
		vendorvalidation.Field(&c.DBClient, vendorvalidation.Required),
	)
}
