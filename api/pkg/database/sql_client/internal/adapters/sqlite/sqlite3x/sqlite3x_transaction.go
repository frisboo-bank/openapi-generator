package sqlite3x

import (
	"context"
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	sqlxutils "frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal/adapters/utils/sqlx"

	"github.com/jmoiron/sqlx"
)

var _ contracts.SQLXTransaction = (*sqlite3SQLXTransaction)(nil)

type sqlite3SQLXTransaction struct {
	tx *sqlx.Tx
}

func NewSQLite3SQLXTransaction(tx *sqlx.Tx) contracts.SQLXTransaction {
	return &sqlite3SQLXTransaction{
		tx: tx,
	}
}

func (s *sqlite3SQLXTransaction) NamedExec(ctx context.Context, query string, args map[string]any) (sql.Result, error) {
	res, err := sqlxutils.NamedExec(ctx, s.tx, query, args)
	if err != nil {
		return nil, fmt.Errorf("NamedExec: %w", err)
	}
	return res, nil
}

func (s *sqlite3SQLXTransaction) NamedGet(ctx context.Context, dest any, query string, args map[string]any) error {
	if err := sqlxutils.NamedGet(ctx, s.tx, dest, query, args); err != nil {
		return fmt.Errorf("NamedGet: %w", err)
	}
	return nil
}

func (s *sqlite3SQLXTransaction) NamedQuery(ctx context.Context, query string, args map[string]any) (contracts.SQLXRows, error) {
	res, err := sqlxutils.NamedQuery(ctx, s.tx, query, args)
	if err != nil {
		return nil, fmt.Errorf("NamedQuery: %w", err)
	}
	return newSQLXRows(res), nil
}

func (s *sqlite3SQLXTransaction) NamedSelect(ctx context.Context, dest any, query string, args map[string]any) error {
	if err := sqlxutils.NamedSelect(ctx, s.tx, dest, query, args); err != nil {
		return fmt.Errorf("NamedSelect: %w", err)
	}
	return nil
}

func (s *sqlite3SQLXTransaction) Commit(ctx context.Context) error   { return s.tx.Commit() }
func (s *sqlite3SQLXTransaction) Rollback(ctx context.Context) error { return s.tx.Rollback() }
