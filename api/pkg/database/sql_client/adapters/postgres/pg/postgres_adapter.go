package pg

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

	_ "github.com/jackc/pgx/v5/stdlib"
)

var _ contracts.SQLClient = (*postgresSQLClientAdapter)(nil)

type postgresSQLClientAdapter struct {
	name   string
	db     *sql.DB
	cfg    *models.SQLClientOptions
	logger loggerContracts.Logger
}

func NewPostgresSQLClientAdapter(
	name string,
	cfg *models.SQLClientOptions,
	logger loggerContracts.Logger,
) (contracts.SQLClient, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	db, err := postgres.ConnectToPostgres(cfg)
	if err != nil {
		logger.Error(
			"postgres connection failed",
			"name", name,
			"error", err,
		)
		return nil, err
	}
	logger.Info("postgres connected", "name", name)

	return &postgresSQLClientAdapter{
		name:   name,
		db:     db.DB,
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (p *postgresSQLClientAdapter) BeginTx(ctx context.Context, opts *sql.TxOptions) (contracts.SQLTransaction, error) {
	tx, err := p.db.BeginTx(ctx, opts)
	if err != nil {
		p.logger.Error(
			"postgres begin transaction failed",
			"name", p.name,
			"error", err,
		)
		return nil, fmt.Errorf("postgres begin transaction: %w", err)
	}

	return NewPostgresSQLTransaction(tx), nil
}

func (p *postgresSQLClientAdapter) Close() error {
	err := p.db.Close()
	if err != nil {
		p.logger.Error(
			"postgres close failed",
			"name", p.name,
			"error", err,
		)
	}
	return err
}

func (p *postgresSQLClientAdapter) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return p.db.ExecContext(ctx, query, args...)
}

func (p *postgresSQLClientAdapter) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

func (p *postgresSQLClientAdapter) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return p.db.QueryContext(ctx, query, args...)
}

func (p *postgresSQLClientAdapter) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return p.db.QueryRowContext(ctx, query, args...)
}

func (p *postgresSQLClientAdapter) Name() string { return p.name }
func (p *postgresSQLClientAdapter) DB() *sql.DB  { return p.db }
func (p *postgresSQLClientAdapter) Type() sqlclienttype.SqlClientType {
	return sqlclienttype.SqlClientTypes.POSTGRES
}
func (p *postgresSQLClientAdapter) Logger() loggerContracts.Logger { return p.logger }
