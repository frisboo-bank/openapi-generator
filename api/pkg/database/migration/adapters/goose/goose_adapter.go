package goose

import (
	"context"
	"database/sql"
	"strconv"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/models"
	migrationtype "frisboo-bank/openapi-generator-service/pkg/database/migration/models/enums/migration_type"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	vGoose "github.com/pressly/goose/v3"
)

var _ contracts.Migration = (*gooseAdapter)(nil)

type gooseAdapter struct {
	name          string
	migrationsDir string
	debug         bool
	db            *sql.DB
	logger        loggerContracts.Logger
}

func NewGooseAdapter(
	name string,
	cfg *models.MigrationOptions,
	database *sql.DB,
	logger loggerContracts.Logger,
) (contracts.Migration, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("database", database)
	validation.AssertNotNil("logger", logger)

	return &gooseAdapter{
		name:          name,
		migrationsDir: cfg.MigrationsDir,
		debug:         cfg.Debug,
		db:            database,
		logger:        logger,
	}, nil
}

func (g *gooseAdapter) Down(ctx context.Context, version uint) error {
	if version == 0 {
		return vGoose.RunContext(ctx, "down", g.db, g.migrationsDir)
	}
	return vGoose.RunContext(ctx, "down-to", g.db, g.migrationsDir, strconv.FormatUint(uint64(version), 10))
}

func (g *gooseAdapter) Up(ctx context.Context, version uint) error {
	if version == 0 {
		return vGoose.RunContext(ctx, "up", g.db, g.migrationsDir)
	}
	return vGoose.RunContext(ctx, "up-to", g.db, g.migrationsDir, strconv.FormatUint(uint64(version), 10))
}

func (g *gooseAdapter) Reset(ctx context.Context) error {
	return vGoose.ResetContext(ctx, g.db, g.migrationsDir)
}

func (g *gooseAdapter) Logger() loggerContracts.Logger {
	return g.logger
}

func (g *gooseAdapter) Name() string {
	return g.name
}

func (g *gooseAdapter) Type() migrationtype.MigrationType {
	return migrationtype.MigrationTypes.GOOSE
}
