package pgx

import (
	"context"
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	sqlxutils "frisboo-bank/openapi-generator-service/pkg/database/sql_client/internal/adapters/utils/sqlx"

	"github.com/jmoiron/sqlx"
)

var _ contracts.SQLXTransaction = (*postgresSQLXTransaction)(nil)

type postgresSQLXTransaction struct {
	tx *sqlx.Tx
}

func NewPostgresSQLXTransaction(tx *sqlx.Tx) contracts.SQLXTransaction {
	return &postgresSQLXTransaction{
		tx: tx,
	}
}

func (p *postgresSQLXTransaction) NamedExec(ctx context.Context, query string, args map[string]any) (sql.Result, error) {
	res, err := sqlxutils.NamedExec(ctx, p.tx, query, args)
	if err != nil {
		return nil, fmt.Errorf("NamedExec: %w", err)
	}
	return res, nil
}

func (p *postgresSQLXTransaction) NamedGet(ctx context.Context, dest any, query string, args map[string]any) error {
	if err := sqlxutils.NamedGet(ctx, p.tx, dest, query, args); err != nil {
		return fmt.Errorf("NamedGet: %w", err)
	}
	return nil
}

func (p *postgresSQLXTransaction) NamedQuery(ctx context.Context, query string, args map[string]any) (contracts.SQLXRows, error) {
	res, err := sqlxutils.NamedQuery(ctx, p.tx, query, args)
	if err != nil {
		return nil, fmt.Errorf("NamedQuery: %w", err)
	}
	return newSQLXRows(res), nil
}

func (p *postgresSQLXTransaction) NamedSelect(ctx context.Context, dest any, query string, args map[string]any) error {
	if err := sqlxutils.NamedSelect(ctx, p.tx, dest, query, args); err != nil {
		return fmt.Errorf("NamedSelect: %w", err)
	}
	return nil
}

func (p *postgresSQLXTransaction) Commit(_ context.Context) error   { return p.tx.Commit() }
func (p *postgresSQLXTransaction) Rollback(_ context.Context) error { return p.tx.Rollback() }
