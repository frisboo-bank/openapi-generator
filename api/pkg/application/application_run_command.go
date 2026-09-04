package application

import (
	"frisboo-bank/openapi-generator-service/pkg/cli"
	cliContracts "frisboo-bank/openapi-generator-service/pkg/cli/contracts"
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"

	"github.com/spf13/cobra"
)

type ApplicationRunCommandOptions struct {
	Short     string
	Long      string
	Bootstrap func(configLoader configContracts.ConfigLoader, env environmentEnum.Environment, cmd *cobra.Command, args []string) error
}

func NewApplicationRunCommand(cfg ApplicationRunCommandOptions) cliContracts.Command {
	return cli.NewCommand(cli.CommandOptions{
		Use:       "run",
		Short:     "Run your application",
		Long:      cfg.Long,
		Bootstrap: cfg.Bootstrap,
	})
}
