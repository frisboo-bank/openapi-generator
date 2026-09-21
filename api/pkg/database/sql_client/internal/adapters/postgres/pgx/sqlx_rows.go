package pgx

import (
	"database/sql"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"

	"github.com/jmoiron/sqlx"
)

var _ contracts.SQLXRows = (*sqlxRows)(nil)

type sqlxRows struct {
	rows *sqlx.Rows
}

func newSQLXRows(rows *sqlx.Rows) contracts.SQLXRows {
	if rows == nil {
		return nil
	}
	return &sqlxRows{rows: rows}
}

func (r *sqlxRows) Close() error                            { return r.rows.Close() }
func (r *sqlxRows) ColumnTypes() ([]*sql.ColumnType, error) { return r.rows.ColumnTypes() }
func (r *sqlxRows) Columns() ([]string, error)              { return r.rows.Columns() }
func (r *sqlxRows) Err() error                              { return r.rows.Err() }
func (r *sqlxRows) Next() bool                              { return r.rows.Next() }
func (r *sqlxRows) NextResultSet() bool                     { return r.rows.NextResultSet() }
func (r *sqlxRows) Scan(dest ...any) error                  { return r.rows.Scan(dest...) }
func (r *sqlxRows) StructScan(dest any) error               { return r.rows.StructScan(dest) }
func (r *sqlxRows) MapScan(dest map[string]any) error       { return r.rows.MapScan(dest) }
func (r *sqlxRows) SliceScan() ([]any, error)               { return r.rows.SliceScan() }
