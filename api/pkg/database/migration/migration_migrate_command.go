package migration

import (
	"context"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/application"
	"frisboo-bank/openapi-generator-service/pkg/cli"
	cliContracts "frisboo-bank/openapi-generator-service/pkg/cli/contracts"
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/container"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	sqlclient "frisboo-bank/openapi-generator-service/pkg/database/sql_client"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"

	"github.com/spf13/cobra"
)

type MigrationMigrateCommandOptions struct {
	Long string
}

func NewMigrationMigrateCommand(cfg MigrationMigrateCommandOptions) cliContracts.Command {
	upCmd := cli.NewCommand(cli.CommandOptions{
		Use:   "up [name]",
		Short: "Upgrade database version",
		Prepare: func(cmd *cobra.Command) {
			cmd.Args = cobra.ExactArgs(1)
			cmd.Flags().Uint("version", 0, "Migration version")
		},
		Bootstrap: func(configLoader configContracts.ConfigLoader, env environmentEnum.Environment, cmd *cobra.Command, args []string) error {
			version, err := cmd.Flags().GetUint("version")
			if err != nil {
				return err
			}

			return executeMigration(configLoader, env, args[0], func(migration contracts.Migration) error {
				return migration.Up(context.Background(), version)
			})
		},
	})

	downCmd := cli.NewCommand(cli.CommandOptions{
		Use:   "down [name]",
		Short: "Downgrade database version",
		Prepare: func(cmd *cobra.Command) {
			cmd.Args = cobra.ExactArgs(1)
			cmd.Flags().Uint("version", 0, "Migration version")
		},
		Bootstrap: func(configLoader configContracts.ConfigLoader, env environmentEnum.Environment, cmd *cobra.Command, args []string) error {
			version, err := cmd.Flags().GetUint("version")
			if err != nil {
				return err
			}

			return executeMigration(configLoader, env, args[0], func(migration contracts.Migration) error {
				return migration.Down(context.Background(), version)
			})
		},
	})

	resetCmd := cli.NewCommand(cli.CommandOptions{
		Use:   "reset [name]",
		Short: "Reset database",
		Prepare: func(cmd *cobra.Command) {
			cmd.Args = cobra.ExactArgs(1)
		},
		Bootstrap: func(configLoader configContracts.ConfigLoader, env environmentEnum.Environment, cmd *cobra.Command, args []string) error {
			return executeMigration(configLoader, env, args[0], func(migration contracts.Migration) error {
				return migration.Reset(context.Background())
			})
		},
	})

	return cli.NewCommand(cli.CommandOptions{
		Use:   "migrate",
		Short: "Run the db migrations",
		Long:  cfg.Long,
		Commands: []cliContracts.Command{
			upCmd,
			downCmd,
			resetCmd,
		},
	})
}

func executeMigration(
	configLoader configContracts.ConfigLoader,
	env environmentEnum.Environment,
	migrationName string,
	cb func(migration contracts.Migration) error,
) error {
	mod := container.NewModule(
		"migration-runner",
		sqlclient.SQLClientModule(env, configLoader),
		MigrationModule(env, configLoader),
	)

	appBuilder := application.NewApplicationBuilder(configLoader, env)
	appBuilder.ProvideModule(mod)

	app, err := appBuilder.Build()
	if err != nil {
		return fmt.Errorf("migration failed with error: %w", err)
	}

	app.ResolveFunc(func(migrations map[string]contracts.Migration) error {
		migration, ok := migrations[migrationName]
		if !ok {
			return fmt.Errorf("migration %q not found", migrationName)
		}

		return cb(migration)
	})

	return nil
}
