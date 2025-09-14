package db

import "database/sql"

type DatabaseInterface interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row

	ExecDebug(query string, args ...any)
	QueryDebug(query string, args ...any)
	Close()
}
