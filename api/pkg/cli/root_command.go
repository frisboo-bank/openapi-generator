package cli

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/cli/contracts"

	"github.com/spf13/cobra"
)

type RootCommandOptions struct {
	Use      string
	Short    string
	Long     string
	Header   func()
	Commands []contracts.Command
}

var _ contracts.RootCommand = (*rootCommand)(nil)

type rootCommand struct {
	rootCmd *cobra.Command
}

func NewRootCommand(cfg *RootCommandOptions) contracts.RootCommand {
	if cfg.Header != nil {
		fmt.Println()
		cfg.Header()
		fmt.Println()
	}

	rootCmd := &cobra.Command{
		Use:   cfg.Use,
		Short: cfg.Short,
	}

	rootCmd.PersistentFlags().String("configPath", "configs", "Path to configuration directory")
	rootCmd.PersistentFlags().String("configName", "application", "Configuration file name (without extension)")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")
	rootCmd.PersistentFlags().String("environment", "development", "Environment (development, testing, preprod, production)")
	rootCmd.PersistentFlags().String("envPrefix", "APP_", "Prefix for environment variables")

	c := &rootCommand{
		rootCmd: rootCmd,
	}

	for _, cmd := range cfg.Commands {
		c.rootCmd.AddCommand(toCobraCommand(cmd))
	}

	return c
}

func (r *rootCommand) Execute() error {
	return r.rootCmd.Execute()
}
