package repositories

import (
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

type SQLXRepository[TModel any] struct {
	sqlClient contracts.SQLXClient
	logger    loggerContracts.Logger
	tableName string
}

func NewSQLXRepository[TModel any](
	sqlClient contracts.SQLXClient,
	logger loggerContracts.Logger,
	tableName string,
) *SQLXRepository[TModel] {
	return &SQLXRepository[TModel]{
		sqlClient: sqlClient,
		logger:    logger,
		tableName: tableName,
	}
}
