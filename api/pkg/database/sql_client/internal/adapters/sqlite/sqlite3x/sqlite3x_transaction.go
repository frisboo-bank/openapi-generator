package sqlite3x

import (
	"context"
	"database/sql"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"

	"github.com/jmoiron/sqlx"
)

var _ contracts.SQLXTransaction = (*sqlite3SQLXTransaction)(nil)

type sqlite3SQLXTransaction struct {
	tx  *sqlx.Tx
	ctx context.Context
}

func NewSQLite3SQLXTransaction(tx *sqlx.Tx, ctx context.Context) contracts.SQLXTransaction {
	return &sqlite3SQLXTransaction{
		tx:  tx,
		ctx: ctx,
	}
}

func (s *sqlite3SQLXTransaction) Commit(ctx context.Context) error {
	panic("unimplemented")
}

func (s *sqlite3SQLXTransaction) NamedExec(ctx context.Context, query string, args map[string]any) (sql.Result, error) {
	panic("unimplemented")
}

func (s *sqlite3SQLXTransaction) NamedGet(ctx context.Context, dest any, query string, args map[string]any) error {
	panic("unimplemented")
}

func (s *sqlite3SQLXTransaction) NamedQuery(ctx context.Context, query string, args map[string]any) (contracts.SQLXRows, error) {
	panic("unimplemented")
}

func (s *sqlite3SQLXTransaction) NamedSelect(ctx context.Context, dest any, query string, args map[string]any) error {
	panic("unimplemented")
}

func (s *sqlite3SQLXTransaction) Rollback(ctx context.Context) error {
	panic("unimplemented")
}
