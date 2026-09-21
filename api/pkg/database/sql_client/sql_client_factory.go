package sqlclient

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/decorators/telemetry/sql"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/decorators/telemetry/sqlx"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal/adapters/postgres/pg"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal/adapters/postgres/pgx"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/models"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	metricscontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/contracts"
	tracercontracts "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

func CreateSQLClient(name string, cfg *models.SQLClientOptions, logger loggerContracts.Logger, tracer tracercontracts.Tracer, metrics metricscontracts.Metrics) (contracts.SQLClientCore, error) {
	validation.AssertNotNil("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)
	validation.AssertNotNil("tracer", tracer)
	validation.AssertNotNil("metrics", metrics)

	switch cfg.Type {
	case sqlclienttype.SqlClientTypes.POSTGRES:
		adapter, err := pg.NewPostgresSQLClientAdapter(name, cfg, logger)
		if err != nil {
			return nil, err
		}
		decoratedAdapter := sql.WrapSQLClientForTelemetry(name, adapter, tracer, metrics)
		return &sqlClient{adapter: decoratedAdapter}, nil

	case sqlclienttype.SqlClientTypes.POSTGRESX:
		adapter, err := pgx.NewPostgresSQLXClientAdapter(name, cfg, logger)
		if err != nil {
			return nil, err
		}
		decoratedAdapter := sqlx.WrapSQLXClientForTelemetry(name, adapter, tracer, metrics)
		return &sqlXClient{adapter: decoratedAdapter}, nil

	default:
		return nil, fmt.Errorf("unsupported SQLClient type: %v", cfg.Type)
	}
}
