package sqlclient

import (
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
	logger loggerContracts.Logger,
	tracer tracercontracts.Tracer,
	metrics metricscontracts.Metrics,
) (contracts.SQLClientCore, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)
	validation.AssertNotNil("tracer", tracer)
	validation.AssertNotNil("metrics", metrics)

	var adapter contracts.SQLXClientAdapter
	var err error

	switch cfg.Type {
	case sqlclienttype.SqlClientTypes.POSTGRESX:
		adapter, err = pgx.NewPostgresSQLXClientAdapter(name, cfg, logger)
	case sqlclienttype.SqlClientTypes.SQLITE3X:
		adapter, err = sqlite3x.NewSQLite3SQLXClientAdapter(name, cfg, logger)
	default:
		err = fmt.Errorf("unsupported SQLClient type: %v", cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	delegate := sqlx.WrapSQLXClientForTelemetry(name, &sqlXClient{adapter: adapter}, tracer, metrics)
	return delegate, nil
}
