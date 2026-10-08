package postgres_test

import (
	"context"
	"testing"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/config"
	sqlclient "frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal"
	sqlclientsslmode "frisboo-bank/openapi-generator-service/pkg/database/sql_client/types/sqlclientsslmode"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/types/sqlclienttype"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/metrics"
	"frisboo-bank/openapi-generator-service/pkg/telemetry/tracer"
	"frisboo-bank/openapi-generator-service/pkg/tests/devcontainers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCreateSQLClient_Postgres(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pg := devcontainers.NewPostgresTestContainer(devcontainers.PostgresTestContainerOptions{
		DBName:                   "sqlclient_test",
		Username:                 "postgres",
		Password:                 "postgres",
		AdditionalWaitStrategies: []wait.Strategy{wait.ForListeningPort("5432/tcp")},
	})
	pg.Run(t)

	host, err := pg.Container().Host(ctx)
	require.NoError(t, err, "failed to get container host")

	port, err := pg.Container().MappedPort(ctx, "5432/tcp")
	require.NoError(t, err, "failed to get container port")

	log := logger.CreateNoopLogger("test", environmentenum.Environments.TESTING)

	tr, err := tracer.CreateNoopTracer("test", log)
	require.NoError(t, err, "failed to create noop tracer")

	me, err := metrics.CreateNoopMetrics("test", log)
	require.NoError(t, err, "failed to create noop metrics")

	client, err := sqlclient.CreateSQLClient("main", &config.SQLClientOptions{
		IsEnabled:     true,
		Type:          sqlclienttype.SqlClientTypes.POSTGRESX,
		Host:          host,
		Port:          port.Port(),
		Database:      "sqlclient_test",
		User:          "postgres",
		Password:      "postgres",
		SSLMode:       sqlclientsslmode.SqlClientSSLModes.DISABLED,
		EnableTracing: false,
		EnableMetrics: false,
	}, log, tr, me)

	require.NoError(t, err, "failed to create sql client")

	require.NoError(t, client.Ping(ctx), "failed to ping sql client")

	assert.Equal(t, sqlclienttype.SqlClientTypes.POSTGRESX, client.Type())

	assert.NoError(t, client.Close(ctx), "failed to close sql client")
}
