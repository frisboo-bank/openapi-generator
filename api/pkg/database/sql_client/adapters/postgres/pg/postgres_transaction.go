package pg

import (
	"context"
	"database/sql"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
)

var _ contracts.SQLTransaction = (*postgresSQLTransaction)(nil)

type postgresSQLTransaction struct {
	tx *sql.Tx
}

func NewPostgresSQLTransaction(
	tx *sql.Tx,
) contracts.SQLTransaction {
	return &postgresSQLTransaction{
		tx: tx,
	}
}

func (p *postgresSQLTransaction) Commit() error {
	return p.tx.Commit()
}

func (p *postgresSQLTransaction) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return p.tx.ExecContext(ctx, query, args...)
}

func (p *postgresSQLTransaction) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return p.tx.QueryContext(ctx, query, args...)
}

func (p *postgresSQLTransaction) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return p.tx.QueryRowContext(ctx, query, args...)
}

func (p *postgresSQLTransaction) Rollback() error {
	return p.tx.Rollback()
}
