package engine

import (
	"database/sql"
	"fmt"
)

// MySQL dbEngine
type SQLite struct {
	db *sql.DB
}

// NewSQLiteEngine returns SQLite engine
func NewSqliteEngine(dbFilePath string) *SQLite {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		panic(err)
	}

	return &SQLite{
		db: db,
	}
}

//SQLite does not have an explicit TRUNCATE TABLE command like other databases.
//Instead, it has added a TRUNCATE optimizer to the DELETE statement.
func (sqlite *SQLite) Truncate(table string) error {
	cmd := fmt.Sprintf("DELETE FROM %s", table)

	_, err := sqlite.db.Exec(cmd)
	return err
}

func (sqlite *SQLite) Close() error {
	return sqlite.db.Close()
}