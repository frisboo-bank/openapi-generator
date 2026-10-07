package migration

import (
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/config"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/internal/adapters/goose"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/types/migrationtype"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

func CreateMigration(
	name string,
	db *sql.DB,
	cfg *config.MigrationOptions,
	env environmentenum.Environment,
	logger loggerContracts.Logger,
) (contracts.Migration, error) {
	var adapter contracts.MigrationAdapter
	var err error

	switch cfg.Type {
	case migrationtype.MigrationTypes.GOOSE:
		adapter, err = goose.NewGooseAdapter(name, cfg, db, logger)
	default:
		err = fmt.Errorf("unsupported Migration type: %v", cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	return &migration{adapter: adapter}, nil
}
