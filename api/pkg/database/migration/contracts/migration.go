package contracts

import (
	"context"

	migrationtype "frisboo-bank/openapi-generator-service/pkg/database/migration/types/migrationtype"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

type (
	Migration interface {
		MigrationAdapter
	}

	MigrationAdapter interface {
		Up(ctx context.Context, version uint) error
		Down(ctx context.Context, version uint) error
		Reset(ctx context.Context) error
		Status(ctx context.Context) error
		CurrentVersion(ctx context.Context) (int64, error)
		Name() string
		Type() migrationtype.MigrationType
		Logger() loggerContracts.Logger
	}
)
