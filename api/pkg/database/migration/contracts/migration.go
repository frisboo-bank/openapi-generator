package contracts

import (
	"context"

	migrationtype "frisboo-bank/openapi-generator-service/pkg/database/migration/models/enums/migration_type"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

type (
	Migration interface {
		Up(ctx context.Context, version uint) error
		Down(ctx context.Context, version uint) error
		Reset(ctx context.Context) error
		Name() string
		Type() migrationtype.MigrationType
		Logger() loggerContracts.Logger
	}
)
