package contracts

import (
	"database/sql"

	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

type (
	SQLRow interface {
		Scan(dest ...any) error
		Err() error
	}

	SQLRows interface {
		Close() error
		ColumnTypes() ([]*sql.ColumnType, error)
		Columns() ([]string, error)
		Err() error
		Next() bool
		NextResultSet() bool
		Scan(dest ...any) error
	}

	SQLClientCore interface {
		Close() error
		Ping() error
		Name() string
		Type() sqlclienttype.SqlClientType
		Logger() loggerContracts.Logger
	}

	SQLXRows interface {
		SQLRows
		StructScan(dest any) error
		MapScan(dest map[string]any) error
		SliceScan() ([]any, error)
	}

	SQLXClient interface {
		SQLClientCore
		SQLXClientAdapter
	}

	SQLXClientAdapter interface {
		SQLClientCore
		BeginTransaction(opts *sql.TxOptions) (SQLXTransaction, error)
		NamedExec(query string, args map[string]any) (sql.Result, error)
		NamedGet(dest any, query string, args map[string]any) error
		NamedQuery(query string, args map[string]any) (SQLXRows, error)
		NamedSelect(dest any, query string, args map[string]any) error
	}

	SQLXTransaction interface {
		Commit() error
		NamedExec(query string, args map[string]any) (sql.Result, error)
		NamedGet(dest any, query string, args map[string]any) error
		NamedQuery(query string, args map[string]any) (SQLXRows, error)
		NamedSelect(dest any, query string, args map[string]any) error
		Rollback() error
	}

	WithDBGetter interface {
		DB() *sql.DB
	}
)
