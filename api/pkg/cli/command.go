package cli

import (
	"fmt"
	"strings"

	"frisboo-bank/openapi-generator-service/pkg/cli/contracts"
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"github.com/spf13/cobra"
)

var _ contracts.Command = (*command)(nil)

type CommandOptions struct {
	Use       string
	Short     string
	Long      string
	Commands  []contracts.Command
	Prepare   func(cmd *cobra.Command)
	Bootstrap func(configLoader configContracts.ConfigLoader, env environmentEnum.Environment, cmd *cobra.Command, args []string) error
}

type command struct {
	use       string
	short     string
	long      string
	commands  []contracts.Command
	prepare   func(cmd *cobra.Command)
	bootstrap func(configLoader configContracts.ConfigLoader, env environmentEnum.Environment, cmd *cobra.Command, args []string) error
}

func NewCommand(cfg CommandOptions) contracts.Command {
	validation.AssertNotEmpty("Use", cfg.Use)

	return &command{
		use:       cfg.Use,
		short:     cfg.Short,
		long:      cfg.Long,
		commands:  cfg.Commands,
		prepare:   cfg.Prepare,
		bootstrap: cfg.Bootstrap,
	}
}

func (c *command) Long() string {
	return c.long
}

func (c *command) Short() string {
	return c.short
}

func (c *command) Use() string {
	return c.use
}

func (c *command) Commands() []contracts.Command {
	return c.commands
}

func (c *command) Run(configLoader configContracts.ConfigLoader, env environmentEnum.Environment, cmd *cobra.Command, args []string) error {
	if len(c.commands) > 0 && c.bootstrap == nil {
		var subCmds []string
		for _, subCmd := range c.commands {
			subCmds = append(subCmds, subCmd.Use())
		}
		return fmt.Errorf("%q requires a subcommand: %q", c.use, strings.Join(subCmds, "|"))
	}
	if c.bootstrap == nil {
		return nil
	}
	return c.bootstrap(configLoader, env, cmd, args)
}

func (c *command) Prepare(cmd *cobra.Command) {
	if c.prepare == nil {
		return
	}
	c.prepare(cmd)
}
