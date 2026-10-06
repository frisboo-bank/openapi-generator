package sqlclient

import (
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

func (s *sqlXClient) NamedExec(query string, args map[string]any) (sql.Result, error) {
	return s.adapter.NamedExec(query, args)
}

func (s *sqlXClient) NamedGet(dest any, query string, args map[string]any) error {
	return s.adapter.NamedGet(dest, query, args)
}

func (s *sqlXClient) NamedQuery(query string, args map[string]any) (contracts.SQLXRows, error) {
	return s.adapter.NamedQuery(query, args)
}

func (s *sqlXClient) NamedSelect(dest any, query string, args map[string]any) error {
	return s.adapter.NamedSelect(dest, query, args)
}

func (s *sqlXClient) BeginTransaction(opts *sql.TxOptions) (contracts.SQLXTransaction, error) {
	return s.adapter.BeginTransaction(opts)
}

func (s *sqlXClient) Close() error { return s.adapter.Close() }

func (s *sqlXClient) DB() *sql.DB {
	if getter, ok := s.adapter.(contracts.WithDBGetter); ok {
		return getter.DB()
	}
	return nil
}

func (s *sqlXClient) Ping() error { return s.adapter.Ping() }

func (s *sqlXClient) Logger() loggerContracts.Logger    { return s.adapter.Logger() }
func (s *sqlXClient) Name() string                      { return s.adapter.Name() }
func (s *sqlXClient) Type() sqlclienttype.SqlClientType { return s.adapter.Type() }
