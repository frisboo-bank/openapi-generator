package contracts

import (
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"

	"github.com/spf13/cobra"
)

type Command interface {
	Use() string
	Short() string
	Long() string
	Commands() []Command
	Prepare(cmd *cobra.Command)
	Run(configLoader configContracts.ConfigLoader, env environmentEnum.Environment, cmd *cobra.Command, args []string) error
}
