package goose_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/database/migration"
	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/tests/devcontainers"
	"frisboo-bank/openapi-generator-service/pkg/utils"

	"github.com/stretchr/testify/require"
)

var (
	migrationDir        = "pkg/database/migration/testdata/migrations/goose"
	latestVersion int64 = 2
)

func createMigrator(t *testing.T) contracts.Migration {
	t.Helper()

	root, err := utils.GetProjectRootWorkingDirectory()
	require.NoError(t, err, "failed to get project root")

	container := devcontainers.NewPostgresTestContainer(devcontainers.PostgresTestContainerOptions{
		DBName: "migration_test",
	})
	container.Run(t)

	migrator, err := migration.CreateMigrationForTests("test", container.DB(t),
		filepath.Join(root, migrationDir),
		environmentenum.Environments.TESTING,
	)
	require.NoError(t, err, "failed to create test migrator")

	return migrator
}

func TestMigrate(t *testing.T) {
	migrator := createMigrator(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	require.NoError(t, migrator.Up(ctx, 0))

	version, err := migrator.CurrentVersion(ctx)
	require.NoError(t, err, "failed to get current version")
	require.Equal(t, latestVersion, version, "migrate up failed")

	require.NoError(t, migrator.Down(ctx, 1))

	version, err = migrator.CurrentVersion(ctx)
	require.NoError(t, err, "failed to get current version")
	require.Equal(t, int64(1), version, "migrate down failed")

	require.NoError(t, migrator.Reset(ctx))

	version, err = migrator.CurrentVersion(ctx)
	require.NoError(t, err, "failed to get current version")
	require.Equal(t, int64(0), version, "migrate reset failed")
}

func TestMigrationCurrentVersionOnPristineDB(t *testing.T) {
	migrator := createMigrator(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	version, err := migrator.CurrentVersion(ctx)
	require.NoError(t, err)
	require.Equal(t, int64(0), version)
}
