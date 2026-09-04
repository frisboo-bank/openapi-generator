package pgx

import (
	"context"
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"

	"github.com/jmoiron/sqlx"
)

var _ contracts.SQLXTransaction = (*postgresSQLXTransaction)(nil)

type postgresSQLXTransaction struct {
	tx *sqlx.Tx
}

func NewPostgresSQLXTransaction(
	tx *sqlx.Tx,
) contracts.SQLXTransaction {
	return &postgresSQLXTransaction{
		tx: tx,
	}
}

func (p *postgresSQLXTransaction) Commit() error {
	return p.tx.Commit()
}

func (p *postgresSQLXTransaction) NamedExec(ctx context.Context, query string, args any) (sql.Result, error) {
	return p.tx.NamedExecContext(ctx, query, args)
}

func (p *postgresSQLXTransaction) NamedGet(ctx context.Context, dest any, query string, args any) error {
	query, args, err := sqlx.Named(query, args)
	if err != nil {
		return fmt.Errorf("named get: %w", err)
	}

	err = p.tx.GetContext(ctx, dest, query, args)
	if err != nil {
		return fmt.Errorf("named get: %w", err)
	}

	return nil
}

func (p *postgresSQLXTransaction) NamedSelect(ctx context.Context, dest any, query string, args any) error {
	query, args, err := sqlx.Named(query, args)
	if err != nil {
		return fmt.Errorf("named select: %w", err)
	}

	err = p.tx.SelectContext(ctx, dest, query, args)
	if err != nil {
		return fmt.Errorf("named select: %w", err)
	}

	return nil
}

func (p *postgresSQLXTransaction) Rollback() error {
	return p.tx.Rollback()
}
