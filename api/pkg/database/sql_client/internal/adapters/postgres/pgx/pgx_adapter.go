package pgx

import (
	"context"
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal/adapters/postgres"
	sqlxutils "frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal/adapters/utils/sqlx"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/config"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/types/sqlclienttype"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	_ contracts.SQLXClientAdapter = (*postgresSQLXClientAdapter)(nil)
	_ contracts.WithDBGetter      = (*postgresSQLXClientAdapter)(nil)
)

type postgresSQLXClientAdapter struct {
	name   string
	cfg    *config.SQLClientOptions
	db     *sqlx.DB
	logger loggerContracts.Logger
}

func NewPostgresSQLXClientAdapter(
	name string,
	cfg *config.SQLClientOptions,
	logger loggerContracts.Logger,
) (contracts.SQLXClientAdapter, error) {
	validation.AssertNotNil("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	db, err := postgres.ConnectToPostgres(cfg)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return &postgresSQLXClientAdapter{
		name:   name,
		cfg:    cfg,
		db:     db,
		logger: logger,
	}, nil
}

func (p *postgresSQLXClientAdapter) BeginTransaction(ctx context.Context, opts *sql.TxOptions) (contracts.SQLXTransaction, error) {
	tx, err := p.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("begin transaction failed with error: %w", err)
	}
	return NewPostgresSQLXTransaction(tx), nil
}

func (p *postgresSQLXClientAdapter) NamedExec(ctx context.Context, query string, args map[string]any) (sql.Result, error) {
	res, err := sqlxutils.NamedExec(ctx, p.db, query, args, sqlx.DOLLAR)
	if err != nil {
		return nil, fmt.Errorf("fetch with NamedExec failed with error: %w", err)
	}
	return res, nil
}

func (p *postgresSQLXClientAdapter) NamedGet(ctx context.Context, dest any, query string, args map[string]any) error {
	if err := sqlxutils.NamedGet(ctx, p.db, dest, query, args, sqlx.DOLLAR); err != nil {
		return fmt.Errorf("fetch with NamedGet failed with error: %w", err)
	}
	return nil
}

func (p *postgresSQLXClientAdapter) NamedQuery(ctx context.Context, query string, args map[string]any) (contracts.SQLXRows, error) {
	res, err := sqlxutils.NamedQuery(ctx, p.db, query, args, sqlx.DOLLAR)
	if err != nil {
		return nil, fmt.Errorf("fetch with NamedQuery failed with error: %w", err)
	}
	return newSQLXRows(res), nil
}

func (p *postgresSQLXClientAdapter) NamedSelect(ctx context.Context, dest any, query string, args map[string]any) error {
	if err := sqlxutils.NamedSelect(ctx, p.db, dest, query, args, sqlx.DOLLAR); err != nil {
		return fmt.Errorf("fetch with NamedSelect failed with error: %w", err)
	}
	return nil
}

func (p *postgresSQLXClientAdapter) Ping(ctx context.Context) error { return p.db.PingContext(ctx) }

func (p *postgresSQLXClientAdapter) Close(ctx context.Context) error { return p.db.Close() }

func (p *postgresSQLXClientAdapter) DB() *sql.DB                    { return p.db.DB }
func (p *postgresSQLXClientAdapter) Logger() loggerContracts.Logger { return p.logger }
func (p *postgresSQLXClientAdapter) Name() string                   { return p.name }
func (p *postgresSQLXClientAdapter) Type() sqlclienttype.SqlClientType {
	return sqlclienttype.SqlClientTypes.POSTGRESX
}
