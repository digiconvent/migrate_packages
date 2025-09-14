package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"

	"github.com/digiconvent/migrate_packages/.test/log"
	"github.com/digiconvent/migrate_packages/.test/utils"
	_ "github.com/mattn/go-sqlite3"
)

func New(uri string) (DatabaseInterface, error) {
	db, err := sql.Open("sqlite3", uri)
	if err != nil {
		return nil, err
	}

	return &SqliteDatabase{
		DB:    db,
		mutex: sync.Mutex{},
	}, nil
}

type SqliteDatabase struct {
	DB    *sql.DB
	mutex sync.Mutex
}

func (s *SqliteDatabase) Exec(query string, args ...any) (sql.Result, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.DB.Exec(query, args...)
}
func (s *SqliteDatabase) ExecDebug(query string, args ...any) {
	result, err := s.DB.Exec(query, args...)

	if err != nil {
		log.Error("Exec failed: " + err.Error())
	}
	log.Info("Exec: " + query)
	log.Info(args)
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Error("Could not get rows affected: " + err.Error())
	}

	log.Info("Rows affected: " + strconv.Itoa(int(rowsAffected)))

	lastInsertId, err := result.LastInsertId()

	if err != nil {
		log.Error("Could not get last insert id: " + err.Error())
	}

	log.Info("Last insert id: " + strconv.Itoa(int(lastInsertId)))
}

func (s *SqliteDatabase) Query(query string, args ...any) (*sql.Rows, error) {
	return s.DB.Query(query, args...)
}

func (s *SqliteDatabase) QueryDebug(query string, args ...any) {
	log.Info("Query: " + query)
	log.Info(args)
	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return
	}

	columns, err := rows.Columns()

	if err != nil {
		return
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	var table = utils.NewTable(columns)
	for rows.Next() {
		err = rows.Scan(valuePtrs...)

		if err != nil {
			return
		}

		table.AddRow(values...)
	}
	count := len(table.Values)
	fmt.Println("This query has " + strconv.Itoa(count) + " rows")
	fmt.Println(table.Render())
}

func (s *SqliteDatabase) QueryRow(query string, args ...any) *sql.Row {
	row := s.DB.QueryRow(query, args...)
	return row
}

func (s *SqliteDatabase) Close() {
	s.DB.Close()
}
