package pgx

import (
	"context"
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/adapters/postgres"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/models"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"github.com/jmoiron/sqlx"
)

var _ contracts.SQLXClient = (*postgresSQLXClientAdapter)(nil)

type postgresSQLXClientAdapter struct {
	name   string
	db     *sqlx.DB
	cfg    *models.SQLClientOptions
	logger loggerContracts.Logger
}

func NewPostgresSQLXClientAdapter(
	name string,
	cfg *models.SQLClientOptions,
	logger loggerContracts.Logger,
) (contracts.SQLXClient, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	db, err := postgres.ConnectToPostgres(cfg)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return &postgresSQLXClientAdapter{
		name:   name,
		db:     db,
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (p *postgresSQLXClientAdapter) BeginTxx(
	ctx context.Context,
	opts *sql.TxOptions,
) (contracts.SQLXTransaction, error) {
	tx, err := p.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	return NewPostgresSQLXTransaction(tx), nil
}

func (p *postgresSQLXClientAdapter) Close() error {
	if err := p.db.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	return nil
}

func (p *postgresSQLXClientAdapter) NamedExec(ctx context.Context, query string, args any) (sql.Result, error) {
	res, err := p.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("named exec: %w", err)
	}
	return res, nil
}

func (p *postgresSQLXClientAdapter) NamedGet(ctx context.Context, dest any, query string, args any) error {
	query, args, err := sqlx.Named(query, args)
	if err != nil {
		return fmt.Errorf("named get: %w", err)
	}
	if err := p.db.GetContext(ctx, dest, query, args); err != nil {
		return fmt.Errorf("named get: %w", err)
	}
	return nil
}

func (p *postgresSQLXClientAdapter) NamedQuery(ctx context.Context, query string, args any) (*sqlx.Rows, error) {
	res, err := p.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("named query: %w", err)
	}
	return res, nil
}

func (p *postgresSQLXClientAdapter) NamedSelect(ctx context.Context, dest any, query string, args any) error {
	query, args, err := sqlx.Named(query, args)
	if err != nil {
		return fmt.Errorf("named select: %w", err)
	}
	if err := p.db.SelectContext(ctx, dest, query, args); err != nil {
		return fmt.Errorf("named select: %w", err)
	}
	return nil
}

func (p *postgresSQLXClientAdapter) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

func (p *postgresSQLXClientAdapter) Name() string { return p.name }
func (p *postgresSQLXClientAdapter) DB() *sql.DB  { return p.db.DB }
func (p *postgresSQLXClientAdapter) Type() sqlclienttype.SqlClientType {
	return sqlclienttype.SqlClientTypes.POSTGRESX
}
func (p *postgresSQLXClientAdapter) Logger() loggerContracts.Logger { return p.logger }
