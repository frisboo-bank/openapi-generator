package migration

import (
	"context"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	migrationtype "frisboo-bank/openapi-generator-service/pkg/database/migration/models/enums/migration_type"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

var _ contracts.Migration = (*migration)(nil)

type migration struct {
	adapter contracts.MigrationAdapter
}

func (m *migration) Down(ctx context.Context, version uint) error {
	return m.adapter.Down(ctx, version)
}

func (m *migration) Reset(ctx context.Context) error {
	return m.adapter.Reset(ctx)
}

func (m *migration) Up(ctx context.Context, version uint) error {
	return m.adapter.Up(ctx, version)
}

func (m *migration) Name() string                      { return m.adapter.Name() }
func (m *migration) Type() migrationtype.MigrationType { return m.adapter.Type() }
func (m *migration) Logger() loggercontracts.Logger {
	return m.adapter.Logger()
}
