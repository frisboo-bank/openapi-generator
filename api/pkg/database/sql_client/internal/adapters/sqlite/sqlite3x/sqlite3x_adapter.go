package sqlite3x

import (
	"context"
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	sqlxutils "frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal/adapters/utils/sqlx"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/config"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/types/sqlclienttype"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var (
	_ contracts.SQLXClientAdapter = (*sqlite3SQLXClientAdapter)(nil)
	_ contracts.WithDBGetter      = (*sqlite3SQLXClientAdapter)(nil)
)

type sqlite3SQLXClientAdapter struct {
	name   string
	cfg    *config.SQLClientOptions
	db     *sqlx.DB
	logger loggerContracts.Logger
}

func NewSQLite3SQLXClientAdapter(
	name string,
	cfg *config.SQLClientOptions,
	logger loggerContracts.Logger,
) (contracts.SQLXClientAdapter, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	db, err := sqlx.ConnectContext(context.Background(), "sqlite", cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return &sqlite3SQLXClientAdapter{
		name:   name,
		cfg:    cfg,
		db:     db,
		logger: logger,
	}, nil
}

func (s *sqlite3SQLXClientAdapter) BeginTransaction(ctx context.Context, opts *sql.TxOptions) (contracts.SQLXTransaction, error) {
	tx, err := s.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("begin transaction failed with error: %w", err)
	}
	return NewSQLite3SQLXTransaction(tx), nil
}

func (s *sqlite3SQLXClientAdapter) NamedExec(ctx context.Context, query string, args map[string]any) (sql.Result, error) {
	res, err := sqlxutils.NamedExec(ctx, s.db, query, args)
	if err != nil {
		return nil, fmt.Errorf("fetch with NamedExec failed with error: %w", err)
	}
	return res, nil
}

func (s *sqlite3SQLXClientAdapter) NamedGet(ctx context.Context, dest any, query string, args map[string]any) error {
	if err := sqlxutils.NamedGet(ctx, s.db, dest, query, args); err != nil {
		return fmt.Errorf("fetch with NamedGet failed with error: %w", err)
	}
	return nil
}

func (s *sqlite3SQLXClientAdapter) NamedQuery(ctx context.Context, query string, args map[string]any) (contracts.SQLXRows, error) {
	res, err := sqlxutils.NamedQuery(ctx, s.db, query, args)
	if err != nil {
		return nil, fmt.Errorf("fetch with NamedQuery failed with error: %w", err)
	}
	return newSQLXRows(res), nil
}

func (s *sqlite3SQLXClientAdapter) NamedSelect(ctx context.Context, dest any, query string, args map[string]any) error {
	if err := sqlxutils.NamedSelect(ctx, s.db, dest, query, args); err != nil {
		return fmt.Errorf("fetch with NamedSelect failed with error: %w", err)
	}
	return nil
}

func (s *sqlite3SQLXClientAdapter) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }

func (s *sqlite3SQLXClientAdapter) Close(ctx context.Context) error { return s.db.Close() }

func (s *sqlite3SQLXClientAdapter) DB() *sql.DB                   { return s.db.DB }
func (s *sqlite3SQLXClientAdapter) Logger() loggerContracts.Logger { return s.logger }
func (s *sqlite3SQLXClientAdapter) Name() string                   { return s.name }
func (s *sqlite3SQLXClientAdapter) Type() sqlclienttype.SqlClientType {
	return sqlclienttype.SqlClientTypes.SQLITE3X
}
