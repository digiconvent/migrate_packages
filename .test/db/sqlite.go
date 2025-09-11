package db

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/digiconvent_clean_pkg/test/log"
	"github.com/digiconvent_clean_pkg/test/utils"
)

var databasePath = ""
var databases = map[string]DatabaseInterface{}

func NewTestSqliteDB(dbName string) DatabaseInterface {
	if databases[dbName] != nil {
		databases[dbName].Close()
		delete(databases, dbName)
	}
	connection, _ := NewSqliteConnection(dbName, true)

	if connection.pkgDir() != "" {
		if _, err := os.Stat(connection.pkgDir()); err == nil {
			err := connection.MigratePackage(false)
			if err != nil {
				connection.DeleteDatabase()
				panic(err.Error() + ": " + connection.pkgDir())
			}
		}
	}

	return connection
}

func NewSqliteDB(dbName string) DatabaseInterface {
	connection, _ := NewSqliteConnection(dbName, false)

	return connection
}

func NewSqliteConnection(dbName string, test bool) (DatabaseInterface, bool) {
	fresh := true
	dbName = strings.ToLower(dbName)
	is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9\.]*$`).MatchString(dbName)
	if !is_alphanumeric {
		panic(fmt.Sprint("Database name must be alphanumeric: ", dbName))
	}

	var dbPath string
	if test {
		dbPath = path.Join(os.TempDir(), "testd9t", "test", dbName)
	} else {
		dbPath = path.Join(os.Getenv(databasePath), dbName)
	}

	if databases[dbName] == nil {
		var db *sql.DB
		var err error
		err = os.MkdirAll(dbPath, 0755)

		if err != nil {
			log.Error("Could not create database directory: " + dbPath)
		}

		dbPath = path.Join(dbPath, "database.db")
		if _, err := os.Stat(dbPath); err == nil {
			log.Success("Loading existing database at " + dbPath)
			fresh = false
		} else {
			log.Warning("Creating database at " + dbPath)
		}

		db, err = sql.Open("sqlite3", dbPath)

		if err != nil {
			log.Error("Could not create/open database: " + dbPath)
			panic(err)
		}

		_, err = db.Exec("PRAGMA foreign_keys = ON;")
		if err != nil {
			log.Error("Could not enable foreign keys")
			panic(err)
		}

		databases[dbName] = &SqliteDatabase{
			DB:   db,
			name: dbName,
			test: test,
		}
	}

	return databases[dbName], fresh
}

type SqliteDatabase struct {
	DB   *sql.DB
	name string
	test bool
}

func ListPackages() []string {
	var packages []string

	for key := range databases {
		packages = append(packages, key)
	}

	return packages
}

func (s *SqliteDatabase) pkgDir() string {
	workingDir, _ := os.Getwd()
	sep := "/testd9t/testd9t"
	segments := strings.Split(workingDir, sep)
	if len(segments) < 2 {
		return ""
	}

	dir := path.Join(segments[0], sep, "backend", "pkg", s.name)
	return dir
}

func (s *SqliteDatabase) MigratePackage(verbose bool) error {
	log.Info("Migrating package " + s.name)
	dbPath := path.Join(s.pkgDir(), "db")

	versions, err := os.ReadDir(dbPath)

	if err != nil {
		return err
	}

	for _, version := range versions {
		if version.IsDir() {
			migrations, err := os.ReadDir(path.Join(dbPath, version.Name()))

			if err != nil {
				return err
			}

			for _, migration := range migrations {
				sql, err := os.ReadFile(path.Join(dbPath, version.Name(), migration.Name()))
				if err != nil {
					return err
				}

				result, err := s.DB.Exec(string(sql))

				if err != nil {
					log.Error("❌ " + s.name + ":" + migration.Name())
					return err
				} else {
					if verbose {
						log.Success("✅ " + s.name + ":" + migration.Name())
					}
				}

				if result == nil {
					log.Error("Migration did not return a result: " + migration.Name() + " on database " + s.name)
					return nil
				}
			}
		}
	}

	return nil
}

func (s *SqliteDatabase) Dir() string {
	if s.test {
		return path.Join("/tmp", "testd9t", "test", s.name)
	}
	return path.Join(os.Getenv(databasePath), s.name)
}

func (s *SqliteDatabase) DeleteDatabase() {
	s.DB.Close()
	var err error
	if s.test {
		err = os.RemoveAll(s.Dir())
	} else {
		err = os.Remove(s.Dir())
	}
	if err != nil {
		log.Error("Could not delete database: " + s.Dir())
		panic(err)
	}

	if _, err := os.Stat(s.Dir()); err == nil {
		log.Error("Database still exists: " + s.Dir())
	}
}

func (s *SqliteDatabase) Exec(query string, args ...any) (sql.Result, error) {
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
