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
	tx  *sqlx.Tx
	ctx context.Context
}

func NewPostgresSQLXTransaction(tx *sqlx.Tx, ctx context.Context) contracts.SQLXTransaction {
	return &postgresSQLXTransaction{
		tx:  tx,
		ctx: ctx,
	}
}

func (p *postgresSQLXTransaction) NamedExec(query string, args map[string]any) (sql.Result, error) {
	res, err := sqlxutils.NamedExec(p.ctx, p.tx, query, args)
	if err != nil {
		return nil, fmt.Errorf("NamedExec: %w", err)
	}
	return res, nil
}

func (p *postgresSQLXTransaction) NamedGet(dest any, query string, args map[string]any) error {
	if err := sqlxutils.NamedGet(p.ctx, p.tx, dest, query, args); err != nil {
		return fmt.Errorf("NamedGet: %w", err)
	}
	return nil
}

func (p *postgresSQLXTransaction) NamedQuery(query string, args map[string]any) (contracts.SQLXRows, error) {
	res, err := sqlxutils.NamedQuery(p.ctx, p.tx, query, args)
	if err != nil {
		return nil, fmt.Errorf("NamedQuery: %w", err)
	}
	return newSQLXRows(res), nil
}

func (p *postgresSQLXTransaction) NamedSelect(dest any, query string, args map[string]any) error {
	if err := sqlxutils.NamedSelect(p.ctx, p.tx, dest, query, args); err != nil {
		return fmt.Errorf("NamedSelect: %w", err)
	}
	return nil
}

func (p *postgresSQLXTransaction) Commit() error   { return p.tx.Commit() }
func (p *postgresSQLXTransaction) Rollback() error { return p.tx.Rollback() }
