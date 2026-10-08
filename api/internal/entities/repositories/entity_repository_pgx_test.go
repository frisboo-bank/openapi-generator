package repositories_test

import (
	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	"frisboo-bank/openapi-generator-service/internal/entities/repositories"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger"
	"frisboo-bank/openapi-generator-service/pkg/tests/devcontainers"

	"context"
	"database/sql"

	sqlclientContracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/types/sqlclienttype"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type EntityRepositoryPgxTestSuite struct {
	suite.Suite
	pg         *devcontainers.PostgresTestContainer
	db         *sqlx.DB
	repository contracts.EntitySQLRepository
}

func (s *EntityRepositoryPgxTestSuite) SetupSuite() {
	s.pg = devcontainers.NewPostgresTestContainer(devcontainers.PostgresTestContainerOptions{})
	s.pg.Run(s.T())

	s.db = s.pg.SQLX(s.T())
	// require.NoError(s.T(), s.runMigrations())

	s.repository = repositories.NewEntityRepositoryPgx(
		&testSQLXAdapter{db: s.db},
		logger.CreateNoopLogger("test", environmentenum.Environments.TESTING),
	)
}

// testSQLXAdapter is a minimal SQLXClientAdapter backed by a *sqlx.DB, used
// by the integration test suite. It exists only because no test-only adapter
// was provided alongside NewEntityRepositoryPgx.
type testSQLXAdapter struct {
	db *sqlx.DB
}

func (t *testSQLXAdapter) Close(ctx context.Context) error                              { return t.db.Close() }
func (t *testSQLXAdapter) Ping(ctx context.Context) error                               { return t.db.PingContext(ctx) }
func (t *testSQLXAdapter) Name() string                                                  { return "test" }
func (t *testSQLXAdapter) Type() sqlclienttype.SqlClientType                             { return sqlclienttype.SqlClientTypes.POSTGRESX }
func (t *testSQLXAdapter) Logger() loggerContracts.Logger                               { return logger.CreateNoopLogger("test", environmentenum.Environments.TESTING) }
func (t *testSQLXAdapter) BeginTransaction(ctx context.Context, opts *sql.TxOptions) (sqlclientContracts.SQLXTransaction, error) {
	tx, err := t.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &testSQLXTransaction{tx}, nil
}
func (t *testSQLXAdapter) NamedExec(ctx context.Context, query string, args map[string]any) (sql.Result, error) {
	return t.db.NamedExecContext(ctx, query, args)
}
func (t *testSQLXAdapter) NamedGet(ctx context.Context, dest any, query string, args map[string]any) error {
	return t.db.GetContext(ctx, dest, query, args)
}
func (t *testSQLXAdapter) NamedQuery(ctx context.Context, query string, args map[string]any) (sqlclientContracts.SQLXRows, error) {
	return t.db.QueryxContext(ctx, query, args)
}
func (t *testSQLXAdapter) NamedSelect(ctx context.Context, dest any, query string, args map[string]any) error {
	return t.db.SelectContext(ctx, dest, query, args)
}

type testSQLXTransaction struct {
	tx *sqlx.Tx
}

func (t *testSQLXTransaction) Commit(ctx context.Context) error   { return t.tx.Commit() }
func (t *testSQLXTransaction) Rollback(ctx context.Context) error  { return t.tx.Rollback() }
func (t *testSQLXTransaction) NamedExec(ctx context.Context, query string, args map[string]any) (sql.Result, error) {
	return t.tx.NamedExecContext(ctx, query, args)
}
func (t *testSQLXTransaction) NamedGet(ctx context.Context, dest any, query string, args map[string]any) error {
	return t.tx.GetContext(ctx, dest, query, args)
}
func (t *testSQLXTransaction) NamedQuery(ctx context.Context, query string, args map[string]any) (sqlclientContracts.SQLXRows, error) {
	return t.tx.QueryxContext(ctx, query, args)
}
func (t *testSQLXTransaction) NamedSelect(ctx context.Context, dest any, query string, args map[string]any) error {
	return t.tx.SelectContext(ctx, dest, query, args)
}
