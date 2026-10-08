package config

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/types/migrationtype"
	"frisboo-bank/openapi-generator-service/pkg/validation/validators"
	vendorvalidation "github.com/go-ozzo/ozzo-validation"
)

var _ configContracts.Configurable = (*MigrationOptions)(nil)

type MigrationOptions struct {
	IsEnabled bool `mapstructure:"enabled"`
	Type      migrationtype.MigrationType `mapstructure:"type"`
	Debug     bool `mapstructure:"debug"`
	MigrationsDir string `mapstructure:"migrationsDir"`

	// EnableTelemetry toggles tracing and metrics for this migration
	// instance. When false (the default) no tracer/metrics are resolved and
	// the adapter runs uninstrumented.
	EnableTelemetry bool `mapstructure:"enableTelemetry"`
	// Tracer / Metrics name the telemetry instances to resolve when
	// EnableTelemetry is true. They are ignored otherwise.
	Tracer  string `mapstructure:"tracer"`
	Metrics string `mapstructure:"metrics"`

	// dependencies
	Logger   string `mapstructure:"logger"`
	DBClient string `mapstructure:"dbClient"`
}

func (c *MigrationOptions) GetEnabled() bool { return c.IsEnabled }
func (c *MigrationOptions) GetLogger() string { return c.Logger }
func (c *MigrationOptions) SetDefaults() {}
func (c *MigrationOptions) Validate() error {
	if !c.IsEnabled {
		return nil
	}
	return vendorvalidation.ValidateStruct(
		vendorvalidation.Field(&c.Type, vendorvalidation.Required, validators.ValidEnum()),
		vendorvalidation.Field(&c.MigrationsDir, vendorvalidation.Required),
		vendorvalidation.Field(&c.DBClient, vendorvalidation.Required),
	)
}
