package main

import (
	"fmt"
	"os"

	"frisboo-bank/openapi-generator-service/internal/shared/app"
	"frisboo-bank/openapi-generator-service/pkg/application"
	"frisboo-bank/openapi-generator-service/pkg/cli"
	cliContracts "frisboo-bank/openapi-generator-service/pkg/cli/contracts"
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/migration"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cli.NewRootCommand(&cli.RootCommandOptions{
		Use:   "openapi-generator-service",
		Short: "Generate OpenAPI schemas for your api",
		Long:  `This OpenAPI Generator service is used by the OpenAPI frontend to generate OpenAPI specifications for your services.`,
		Header: func() {
			ptermLogo, _ := pterm.DefaultBigText.WithLetters(
				putils.LettersFromStringWithStyle("OpenAPI Generator", pterm.NewStyle(pterm.FgLightCyan)),
				putils.LettersFromStringWithStyle(" Service", pterm.NewStyle(pterm.FgLightMagenta)),
			).Srender()
			pterm.DefaultCenter.Print(ptermLogo)
		},
		Commands: []cliContracts.Command{
			application.NewApplicationRunCommand(application.ApplicationRunCommandOptions{
				Bootstrap: func(configLoader configContracts.ConfigLoader, env environmentEnum.Environment, cmd *cobra.Command, args []string) error {
					bootstrap := app.NewBootstrap(configLoader, env)
					return bootstrap.Run()
				},
			}),
			migration.NewMigrationMigrateCommand(migration.MigrationMigrateCommandOptions{}),
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %v\n", err)
		os.Exit(1)
	}
}
