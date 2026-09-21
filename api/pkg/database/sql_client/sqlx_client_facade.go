package sqlclient

import (
	"context"
	"database/sql"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

var (
	_ contracts.SQLClientCore = (*sqlXClient)(nil)
	_ contracts.SQLXClient    = (*sqlXClient)(nil)
	_ contracts.WithDBGetter  = (*sqlXClient)(nil)
)

type sqlXClient struct {
	adapter contracts.SQLXClientAdapter
}

func (s *sqlXClient) NamedExec(ctx context.Context, query string, args map[string]any) (sql.Result, error) {
	return s.adapter.NamedExec(ctx, query, args)
}

func (s *sqlXClient) NamedGet(ctx context.Context, dest any, query string, args map[string]any) error {
	return s.adapter.NamedGet(ctx, dest, query, args)
}

func (s *sqlXClient) NamedQuery(ctx context.Context, query string, args map[string]any) (contracts.SQLXRows, error) {
	return s.adapter.NamedQuery(ctx, query, args)
}

func (s *sqlXClient) NamedSelect(ctx context.Context, dest any, query string, args map[string]any) error {
	return s.adapter.NamedSelect(ctx, dest, query, args)
}

func (s *sqlXClient) BeginTransaction(ctx context.Context, opts *sql.TxOptions) (contracts.SQLXTransaction, error) {
	return s.adapter.BeginTransaction(ctx, opts)
}

func (s *sqlXClient) DB() *sql.DB {
	if getter, ok := s.adapter.(contracts.WithDBGetter); ok {
		return getter.DB()
	}
	return nil
}

func (s *sqlXClient) Close(ctx context.Context) error   { return s.adapter.Close(ctx) }
func (s *sqlXClient) Logger() loggerContracts.Logger    { return s.adapter.Logger() }
func (s *sqlXClient) Name() string                      { return s.adapter.Name() }
func (s *sqlXClient) Ping(ctx context.Context) error    { return s.adapter.Ping(ctx) }
func (s *sqlXClient) Type() sqlclienttype.SqlClientType { return s.adapter.Type() }
