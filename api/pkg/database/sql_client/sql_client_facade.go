package sqlclient

import (
	"context"
	"database/sql"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

var (
	_ contracts.SQLClientCore = (*sqlClient)(nil)
	_ contracts.SQLClient     = (*sqlClient)(nil)
	_ contracts.WithDBGetter  = (*sqlClient)(nil)
)

type sqlClient struct {
	adapter contracts.SQLClientAdapter
}

func (c *sqlClient) BeginTx(ctx context.Context, opts *sql.TxOptions) (contracts.SQLTransaction, error) {
	return c.adapter.BeginTx(ctx, opts)
}

func (c *sqlClient) Close(ctx context.Context) error { return c.adapter.Close(ctx) }

func (c *sqlClient) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return c.adapter.Exec(ctx, query, args...)
}

func (c *sqlClient) Ping(ctx context.Context) error { return c.adapter.Ping(ctx) }

func (c *sqlClient) Query(ctx context.Context, query string, args ...any) (contracts.SQLRows, error) {
	return c.adapter.Query(ctx, query, args...)
}

func (c *sqlClient) QueryRow(ctx context.Context, query string, args ...any) contracts.SQLRow {
	return c.adapter.QueryRow(ctx, query, args...)
}

func (c *sqlClient) DB() *sql.DB {
	if getter, ok := c.adapter.(contracts.WithDBGetter); ok {
		return getter.DB()
	}
	return nil
}

func (c *sqlClient) Name() string                      { return c.adapter.Name() }
func (c *sqlClient) Type() sqlclienttype.SqlClientType { return c.adapter.Type() }
func (c *sqlClient) Logger() loggerContracts.Logger    { return c.adapter.Logger() }
