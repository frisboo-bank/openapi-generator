package sqlite3x

import (
	"context"
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	sqlxutils "frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal/adapters/utils/sqlx"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/models"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"
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
	cfg    *models.SQLClientOptions
	db     *sqlx.DB
	ctx    context.Context
	logger loggerContracts.Logger
}

func NewSQLite3SQLXClientAdapter(
	name string,
	cfg *models.SQLClientOptions,
	ctx context.Context,
	logger loggerContracts.Logger,
) (contracts.SQLXClientAdapter, error) {
	validation.AssertNotNil("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("ctx", ctx)
	validation.AssertNotNil("logger", logger)

	db, err := sqlx.ConnectContext(ctx, "sqlite", cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return sqlite3SQLXClientAdapter{
		name:   name,
		cfg:    cfg,
		ctx:    ctx,
		db:     db,
		logger: logger,
	}, nil
}

func (s sqlite3SQLXClientAdapter) BeginTransaction(opts *sql.TxOptions) (contracts.SQLXTransaction, error) {
	tx, err := s.db.BeginTxx(s.ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("begin transaction failed with error: %w", err)
	}
	return NewSQLite3SQLXTransaction(tx, s.ctx), nil
}

func (s sqlite3SQLXClientAdapter) NamedExec(query string, args map[string]any) (sql.Result, error) {
	panic("unimplemented")
}

func (s sqlite3SQLXClientAdapter) NamedGet(dest any, query string, args map[string]any) error {
	panic("unimplemented")
}

func (s sqlite3SQLXClientAdapter) NamedQuery(query string, args map[string]any) (contracts.SQLXRows, error) {
	panic("unimplemented")
}

func (s sqlite3SQLXClientAdapter) NamedSelect(dest any, query string, args map[string]any) error {
	if err := sqlxutils.NamedSelect(s.ctx, s.db, dest, query, args); err != nil {
		return fmt.Errorf("fetch with NamedSelect failed with error: %w", err)
	}
	return nil
}

func (s sqlite3SQLXClientAdapter) Ping() error { return s.db.PingContext(s.ctx) }

func (s sqlite3SQLXClientAdapter) Close() error { return s.db.Close() }

func (s *sqlite3SQLXClientAdapter) DB() *sql.DB                   { return s.db.DB }
func (s sqlite3SQLXClientAdapter) Logger() loggerContracts.Logger { return s.logger }
func (s sqlite3SQLXClientAdapter) Name() string                   { return s.name }
func (s sqlite3SQLXClientAdapter) Type() sqlclienttype.SqlClientType {
	return sqlclienttype.SqlClientTypes.SQLITE3X
}
