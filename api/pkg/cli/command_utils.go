package cli

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/cli/contracts"
	"frisboo-bank/openapi-generator-service/pkg/config"
	environmentenums "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func toCobraCommand(cmd contracts.Command) *cobra.Command {
	use := cmd.Use()
	short := cmd.Short()

	validation.AssertNotEmpty("use", use)

	newCmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			configPath, err := cobraCmd.Flags().GetString("configPath")
			if err != nil {
				return fmt.Errorf("get configPath: %w", err)
			}
			configName, err := cobraCmd.Flags().GetString("configName")
			if err != nil {
				return fmt.Errorf("get configName: %w", err)
			}
			debug, err := cobraCmd.Flags().GetBool("debug")
			if err != nil {
				return fmt.Errorf("get debug: %w", err)
			}
			envStr, err := cobraCmd.Flags().GetString("environment")
			if err != nil {
				return fmt.Errorf("get environment: %w", err)
			}
			envPrefix, err := cobraCmd.Flags().GetString("envPrefix")
			if err != nil {
				return fmt.Errorf("get envPrefix: %w", err)
			}
			env, err := environmentenums.ParseEnvironment(envStr)
			if err != nil {
				return fmt.Errorf("parse environment: %w", err)
			}

			cfgLoader, err := config.NewConfigLoader(config.ConfigLoaderOptions{
				ConfigName:     configName,
				ConfigPath:     configPath,
				Debug:          debug,
				EnvKeyReplacer: map[string]string{},
				EnvPrefix:      envPrefix,
			}, viper.New())
			if err != nil {
				return fmt.Errorf("command run failed with error: %w", err)
			}

			return cmd.Run(cfgLoader, env, cobraCmd, args)
		},
	}

	for _, childCmd := range cmd.Commands() {
		newCmd.AddCommand(toCobraCommand(childCmd))
	}

	cmd.Prepare(newCmd)

	return newCmd
}
