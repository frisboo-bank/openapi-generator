package sqlclient

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/types/sqlclienttype"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/config"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSQLClient_SQLite3(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Temp file (not :memory:) so the sqlx pool shares one DB across connections.
	dbPath := filepath.Join(t.TempDir(), "sqlclient_test.db")

	log := logger.CreateNoopLogger("test", environmentenum.Environments.TESTING)

	tr, err := tracer.CreateNoopTracer("test", log)
	require.NoError(t, err, "failed to create noop tracer")

	me, err := metrics.CreateNoopMetrics("test", log)
	require.NoError(t, err, "failed to create noop metrics")

	client, err := CreateSQLClient("main", &config.SQLClientOptions{
		IsEnabled:     true,
		Type:          sqlclienttype.SqlClientTypes.SQLITE3X,
		Database:      dbPath,
		EnableTracing: false,
		EnableMetrics: false,
	}, log, tr, me)

	require.NoError(t, err, "failed to create sql client")

	require.NoError(t, client.Ping(ctx), "failed to ping sql client")

	assert.Equal(t, sqlclienttype.SqlClientTypes.SQLITE3X, client.Type())

	assert.NoError(t, client.Close(ctx), "failed to close sql client")
}
