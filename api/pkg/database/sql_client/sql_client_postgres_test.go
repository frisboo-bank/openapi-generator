package sqlclient_test

import (
	"context"
	"testing"
	"time"

	sqlclient "frisboo-bank/openapi-generator-service/pkg/database/sql_client"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/models"
	sqlclientsslmode "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_ssl_mode"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger"
	"frisboo-bank/openapi-generator-service/pkg/tests/devcontainers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCreateSQLClient_Postgres(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	client, err := sqlclient.CreateSQLClient("main", &models.SQLClientOptions{
		IsEnabled:     false,
		Type:          sqlclienttype.SqlClientType{},
		Debug:         false,
		Host:          host,
		Port:          port.Port(),
		Database:      "sqlclient_test",
		User:          "postgres",
		Password:      "postgres",
		SSLMode:       sqlclientsslmode.SqlClientSSLModes.DISABLED,
		EnableTracing: false,
	}, logger.CreateNoopLogger("test", environmentenum.Environments.TESTING), nil, nil)

	require.NoError(t, err, "failed to create sql client")

	require.NoError(t, client.Ping(ctx), "failed to ping sql client")

	assert.Equal(t, sqlclienttype.SqlClientTypes.POSTGRESX, client.Type())

	assert.NoError(t, client.Close(ctx), "failed to close sql client")
}
