package contracts

import (
	"context"
	"database/sql"

	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"github.com/jmoiron/sqlx"
)

type (
	SQLClientCore interface {
		// Close gracefully shuts down the database connection pool.
		Close() error
		// Name returns the logical name assigned to this client (e.g., "main").
		Name() string
		// Type returns the adapter type
		Type() sqlclienttype.SqlClientType
		// Ping verifies that the database is reachable.
		Ping(ctx context.Context) error
		// Logger returns the logger.
		Logger() loggerContracts.Logger
	}

	SQLClient interface {
		SQLClientCore

		// BeginTx starts a new database transaction.
		BeginTx(ctx context.Context, opts *sql.TxOptions) (SQLTransaction, error)
		// Exec executes a query without returning any rows.
		Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
		// Query executes a query that returns rows.
		Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		// QueryRow executes a query that is expected to return at most one row.
		QueryRow(ctx context.Context, query string, args ...any) *sql.Row
	}

	SQLTransaction interface {
		// Exec executes a query inside the transaction.
		Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
		// Query executes a query that returns rows inside the transaction.
		Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		// QueryRow executes a query expected to return at most one row inside the transaction.
		QueryRow(ctx context.Context, query string, args ...any) *sql.Row
		// Commit commits the transaction.
		Commit() error
		// Rollback aborts the transaction.
		Rollback() error
	}

	SQLXClient interface {
		SQLClientCore

		// BeginTxx starts a new sqlx transaction with the given options.
		BeginTxx(ctx context.Context, opts *sql.TxOptions) (SQLXTransaction, error)
		// Get executes a query and scans the first row into dest.
		NamedGet(ctx context.Context, dest any, query string, args any) error
		// NamedExec executes a named‑parameter query without returning rows.
		NamedExec(ctx context.Context, query string, args any) (sql.Result, error)
		// NamedQuery executes a named‑parameter query and returns the rows for scanning.
		NamedQuery(ctx context.Context, query string, args any) (*sqlx.Rows, error)
		// Select executes a query and scans all rows into dest (a slice pointer).
		NamedSelect(ctx context.Context, dest any, query string, args any) error
	}

	SQLXTransaction interface {
		// Get executes a query and scans the first row into dest inside the transaction.
		NamedGet(ctx context.Context, dest any, query string, args any) error
		// NamedExec executes a named‑parameter query without returning rows inside the transaction.
		NamedExec(ctx context.Context, query string, args any) (sql.Result, error)
		// Select executes a query and scans all rows into dest inside the transaction.
		NamedSelect(ctx context.Context, dest any, query string, args any) error
		// Commit commits the transaction.
		Commit() error
		// Rollback aborts the transaction.
		Rollback() error
	}

	WithDBGetter interface {
		DB() *sql.DB
	}
)
