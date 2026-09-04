package cli

import (
	"fmt"

	cacheEnums "frisboo-bank/openapi-generator-service/pkg/cache/models/enums"
	"frisboo-bank/openapi-generator-service/pkg/cli/contracts"
	"frisboo-bank/openapi-generator-service/pkg/config"
	migrationEnums "frisboo-bank/openapi-generator-service/pkg/database/migration/models/enums"
	sqlclientEnums "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums"
	environmentEnums "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	httpServerEnums "frisboo-bank/openapi-generator-service/pkg/http/http_server/models/enums"
	loggerEnums "frisboo-bank/openapi-generator-service/pkg/logger/models/enums"
	paginationEnums "frisboo-bank/openapi-generator-service/pkg/query/models/enums"
	rpcServerEnums "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/models/enums"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"github.com/go-viper/mapstructure/v2"
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
			env, err := environmentEnums.ParseEnvironment(envStr)
			if err != nil {
				return fmt.Errorf("parse environment: %w", err)
			}

			cfgLoader, err := config.NewConfigLoader(config.ConfigLoaderOptions{
				ConfigName:     configName,
				ConfigPath:     configPath,
				Debug:          debug,
				EnvKeyReplacer: map[string]string{},
				EnvPrefix:      envPrefix,
				DecodeHookFuncs: []mapstructure.DecodeHookFunc{
					cacheEnums.CacheEnumsDecodeHook(),
					environmentEnums.EnvironmentEnumsDecodeHook(),
					httpServerEnums.HTTPServerEnumsDecodeHook(),
					loggerEnums.LoggerEnumsDecodeHook(),
					migrationEnums.MigrationEnumsDecodeHook(),
					paginationEnums.PaginationEnumsDecodeHook(),
					rpcServerEnums.RPCServerEnumsDecodeHook(),
					sqlclientEnums.SQLClientEnumsDecodeHook(),
				},
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
