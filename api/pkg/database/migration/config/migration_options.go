package config

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	migrationtype "frisboo-bank/openapi-generator-service/pkg/database/migration/types/migrationtype"
	"frisboo-bank/openapi-generator-service/pkg/validation/validators"

	vendorvalidation "github.com/go-ozzo/ozzo-validation"
)

var _ configContracts.Configurable = (*MigrationOptions)(nil)

type MigrationOptions struct {
	IsEnabled bool `mapstructure:"enabled" json:"enabled"`
	Type migrationtype.MigrationType `mapstructure:"type" json:"type"`
	Debug bool `mapstructure:"debug" json:"debug"`
	MigrationsDir string `mapstructure:"migrationsDir" json:"migrationsDir"`

	// dependencies
	Logger string `mapstructure:"logger" json:"logger"`
	DBClient string `mapstructure:"dbClient" json:"dbClient"`
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

//go:generate go run github.com/invopop/jsonschema -o schema.json -package config MigrationOptions
