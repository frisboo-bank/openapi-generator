package pg

import (
	"database/sql"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
)

var (
	_ contracts.SQLRows = (*sqlRows)(nil)
	_ contracts.SQLRow  = (*sqlRow)(nil)
)

type sqlRows struct {
	rows *sql.Rows
}

func newSQLRows(rows *sql.Rows) contracts.SQLRows {
	if rows == nil {
		return nil
	}
	return &sqlRows{rows: rows}
}

func (r *sqlRows) Close() error                            { return r.rows.Close() }
func (r *sqlRows) ColumnTypes() ([]*sql.ColumnType, error) { return r.rows.ColumnTypes() }
func (r *sqlRows) Columns() ([]string, error)              { return r.rows.Columns() }
func (r *sqlRows) Err() error                              { return r.rows.Err() }
func (r *sqlRows) Next() bool                              { return r.rows.Next() }
func (r *sqlRows) NextResultSet() bool                     { return r.rows.NextResultSet() }
func (r *sqlRows) Scan(dest ...any) error                  { return r.rows.Scan(dest...) }

type sqlRow struct {
	row *sql.Row
}

func newSQLRow(row *sql.Row) contracts.SQLRow {
	if row == nil {
		return nil
	}
	return &sqlRow{row: row}
}

func (r *sqlRow) Scan(dest ...any) error { return r.row.Scan(dest...) }
func (r *sqlRow) Err() error             { return r.row.Err() }
