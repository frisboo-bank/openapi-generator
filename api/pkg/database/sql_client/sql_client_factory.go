package sqlclient

import (
	"context"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/decorators/telemetry/sqlx"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal/adapters/postgres/pgx"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal/adapters/sqlite/sqlite3x"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/models"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

func CreateSQLClient(
	name string,
	cfg *models.SQLClientOptions,
	ctx context.Context,
	logger loggerContracts.Logger,
	tracer tracercontracts.Tracer,
	metrics metricscontracts.Metrics,
) (contracts.SQLClientCore, error) {
	validation.AssertNotNil("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("ctx", ctx)
	validation.AssertNotNil("logger", logger)
	validation.AssertNotNil("tracer", tracer)
	validation.AssertNotNil("metrics", metrics)

	var adapter contracts.SQLClientCore
	var err error

	switch cfg.Type {
	case sqlclienttype.SqlClientTypes.POSTGRESX:
		adapter, err = pgx.NewPostgresSQLXClientAdapter(name, cfg, ctx, logger)
	case sqlclienttype.SqlClientTypes.SQLITE3X:
		adapter, err = sqlite3x.NewSQLite3SQLXClientAdapter(name, cfg, ctx, logger)
	default:
		err = fmt.Errorf("unsupported SQLClient type: %v", cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	if delegate, ok := adapter.(contracts.SQLXClientAdapter); ok {
		delegate = sqlx.WrapSQLXClientForTelemetry(name, delegate, tracer, metrics)
		return &sqlXClient{adapter: delegate}, nil
	}
	return nil, fmt.Errorf("unknown SQLClient type: %v", cfg.Type)
}
