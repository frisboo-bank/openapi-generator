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

// BindNamed converts a named-placeholder query into a driver-specific
// positional query. The placeholder style must match the underlying driver:
// sqlx.DOLLAR for postgres/pgx, sqlx.QUESTION for sqlite3x.
func BindNamed(query string, args map[string]any, style sqlx.Placeholder) (string, []any, error) {
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

	return sqlx.Rebind(style, namedQuery), namedArgs, nil
}

func NamedExec(ctx context.Context, ex SQLXExecutor, query string, args map[string]any, style sqlx.Placeholder) (sql.Result, error) {
	q, a, err := BindNamed(query, args, style)
	if err != nil {
		return nil, err
	}
	return ex.ExecContext(ctx, q, a...)
}

func NamedGet(ctx context.Context, ex SQLXExecutor, dest any, query string, args map[string]any, style sqlx.Placeholder) error {
	q, a, err := BindNamed(query, args, style)
	if err != nil {
		return err
	}
	return ex.GetContext(ctx, dest, q, a...)
}

func NamedQuery(ctx context.Context, ex SQLXExecutor, query string, args map[string]any, style sqlx.Placeholder) (*sqlx.Rows, error) {
	q, a, err := BindNamed(query, args, style)
	if err != nil {
		return nil, err
	}
	return ex.QueryxContext(ctx, q, a...)
}

func NamedSelect(ctx context.Context, ex SQLXExecutor, dest any, query string, args map[string]any, style sqlx.Placeholder) error {
	q, a, err := BindNamed(query, args, style)
	if err != nil {
		return err
	}
	return ex.SelectContext(ctx, dest, q, a...)
}
