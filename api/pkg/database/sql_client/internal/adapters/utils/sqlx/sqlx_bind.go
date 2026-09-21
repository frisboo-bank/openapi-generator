package sqlx

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SQLXExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	GetContext(context.Context, any, string, ...any) error
	QueryxContext(context.Context, string, ...any) (*sqlx.Rows, error)
	SelectContext(context.Context, any, string, ...any) error
}

func BindNamed(query string, args map[string]any) (string, []any, error) {
	if len(args) == 0 {
		return query, nil, nil
	}

	namedQuery, namedArgs, err := sqlx.Named(query, args)
	if err != nil {
		return "", nil, err
	}

	if len(namedArgs) == 0 {
		return "", nil, fmt.Errorf("query has no named placeholders for %d provided args", len(args))
	}

	return sqlx.Rebind(sqlx.DOLLAR, namedQuery), namedArgs, nil
}

func NamedExec(ctx context.Context, ex SQLXExecutor, query string, args map[string]any) (sql.Result, error) {
	q, a, err := BindNamed(query, args)
	if err != nil {
		return nil, err
	}
	return ex.ExecContext(ctx, q, a...)
}

func NamedGet(ctx context.Context, ex SQLXExecutor, dest any, query string, args map[string]any) error {
	q, a, err := BindNamed(query, args)
	if err != nil {
		return err
	}
	return ex.GetContext(ctx, dest, q, a...)
}

func NamedQuery(ctx context.Context, ex SQLXExecutor, query string, args map[string]any) (*sqlx.Rows, error) {
	q, a, err := BindNamed(query, args)
	if err != nil {
		return nil, err
	}
	return ex.QueryxContext(ctx, q, a...)
}

func NamedSelect(ctx context.Context, ex SQLXExecutor, dest any, query string, args map[string]any) error {
	q, a, err := BindNamed(query, args)
	if err != nil {
		return err
	}
	return ex.SelectContext(ctx, dest, q, a...)
}
