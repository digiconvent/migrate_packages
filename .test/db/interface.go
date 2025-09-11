package db

import "database/sql"

type DatabaseInterface interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecDebug(query string, args ...any)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryDebug(query string, args ...any)
	QueryRow(query string, args ...any) *sql.Row
	pkgDir() string
	MigratePackage(verbose bool) error
	DeleteDatabase()
	Close()
}
