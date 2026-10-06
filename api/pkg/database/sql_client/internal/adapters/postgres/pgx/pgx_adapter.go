package pgx

import (
	"context"
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal/adapters/postgres"
	sqlxutils "frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal/adapters/utils/sqlx"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/models"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"github.com/jmoiron/sqlx"
)

var (
	_ contracts.SQLXClientAdapter = (*postgresSQLXClientAdapter)(nil)
	_ contracts.WithDBGetter      = (*postgresSQLXClientAdapter)(nil)
)

type postgresSQLXClientAdapter struct {
	name   string
	cfg    *models.SQLClientOptions
	ctx    context.Context
	db     *sqlx.DB
	logger loggerContracts.Logger
}

func NewPostgresSQLXClientAdapter(
	name string,
	cfg *models.SQLClientOptions,
	ctx context.Context,
	logger loggerContracts.Logger,
) (contracts.SQLXClientAdapter, error) {
	validation.AssertNotNil("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("ctx", ctx)
	validation.AssertNotNil("logger", logger)

	db, err := postgres.ConnectToPostgres(cfg)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return &postgresSQLXClientAdapter{
		name:   name,
		cfg:    cfg,
		ctx:    ctx,
		db:     db,
		logger: logger,
	}, nil
}

func (p *postgresSQLXClientAdapter) BeginTransaction(opts *sql.TxOptions) (contracts.SQLXTransaction, error) {
	tx, err := p.db.BeginTxx(p.ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("begin transaction failed with error: %w", err)
	}
	return NewPostgresSQLXTransaction(tx), nil
}

func (p *postgresSQLXClientAdapter) NamedExec(query string, args map[string]any) (sql.Result, error) {
	res, err := sqlxutils.NamedExec(p.ctx, p.db, query, args)
	if err != nil {
		return nil, fmt.Errorf("fetch with NamedExec failed with error: %w", err)
	}
	return res, nil
}

func (p *postgresSQLXClientAdapter) NamedGet(dest any, query string, args map[string]any) error {
	if err := sqlxutils.NamedGet(p.ctx, p.db, dest, query, args); err != nil {
		return fmt.Errorf("fetch with NamedGet failed with error: %w", err)
	}
	return nil
}

func (p *postgresSQLXClientAdapter) NamedQuery(query string, args map[string]any) (contracts.SQLXRows, error) {
	res, err := sqlxutils.NamedQuery(p.ctx, p.db, query, args)
	if err != nil {
		return nil, fmt.Errorf("fetch with NamedQuery failed with error: %w", err)
	}
	return newSQLXRows(res), nil
}

func (p *postgresSQLXClientAdapter) NamedSelect(dest any, query string, args map[string]any) error {
	if err := sqlxutils.NamedSelect(p.ctx, p.db, dest, query, args); err != nil {
		return fmt.Errorf("fetch with NamedSelect failed with error: %w", err)
	}
	return nil
}

func (p *postgresSQLXClientAdapter) Ping() error { return p.db.PingContext(p.ctx) }

func (p *postgresSQLXClientAdapter) Close() error { return p.db.Close() }

func (p *postgresSQLXClientAdapter) DB() *sql.DB                    { return p.db.DB }
func (p *postgresSQLXClientAdapter) Logger() loggerContracts.Logger { return p.logger }
func (p *postgresSQLXClientAdapter) Name() string                   { return p.name }
func (p *postgresSQLXClientAdapter) Type() sqlclienttype.SqlClientType {
	return sqlclienttype.SqlClientTypes.POSTGRESX
}
