package postgres

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

func BindAndExecute(ctx context.Context, query string, args map[string]any, exec func(query string, args ...any) error) error {
	namedQuery, namedArgs, err := BindNamed(query, args)
	if err != nil {
		return err
	}

	if len(namedArgs) == 0 {
		return exec(namedQuery)
	}
	return exec(namedQuery, namedArgs...)
}
