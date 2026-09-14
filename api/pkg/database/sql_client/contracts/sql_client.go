package contracts

import (
	"context"
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
		Close(ctx context.Context) error
		Ping(ctx context.Context) error
		Name() string
		Type() sqlclienttype.SqlClientType
		Logger() loggerContracts.Logger
	}

	SQLClient interface {
		SQLClientCore
		SQLClientAdapter
	}

	SQLClientAdapter interface {
		SQLClientCore
		BeginTx(ctx context.Context, opts *sql.TxOptions) (SQLTransaction, error)
		Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
		Query(ctx context.Context, query string, args ...any) (SQLRows, error)
		QueryRow(ctx context.Context, query string, args ...any) SQLRow
	}

	SQLTransaction interface {
		Commit(ctx context.Context) error
		Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
		Query(ctx context.Context, query string, args ...any) (SQLRows, error)
		QueryRow(ctx context.Context, query string, args ...any) SQLRow
		Rollback(ctx context.Context) error
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
		BeginTransaction(ctx context.Context, opts *sql.TxOptions) (SQLXTransaction, error)
		NamedExec(ctx context.Context, query string, args map[string]any) (sql.Result, error)
		NamedGet(ctx context.Context, dest any, query string, args map[string]any) error
		NamedQuery(ctx context.Context, query string, args map[string]any) (SQLXRows, error)
		NamedSelect(ctx context.Context, dest any, query string, args map[string]any) error
	}

	SQLXTransaction interface {
		Commit(ctx context.Context) error
		NamedExec(ctx context.Context, query string, args map[string]any) (sql.Result, error)
		NamedGet(ctx context.Context, dest any, query string, args map[string]any) error
		NamedQuery(ctx context.Context, query string, args map[string]any) (SQLXRows, error)
		NamedSelect(ctx context.Context, dest any, query string, args map[string]any) error
		Rollback(ctx context.Context) error
	}

	WithDBGetter interface {
		DB() *sql.DB
	}
)
