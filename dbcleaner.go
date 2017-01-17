package dbcleaner

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

type dbcleaner struct{}

func New(driver, connectionString string) (*dbcleaner, error) {
	var err error
	db, err = sql.Open(driver, connectionString)

	if err != nil {
		return nil, err
	}

	return &dbcleaner{}, err
}

func (c *dbcleaner) Close() error {
	return db.Close()
}

func (c *dbcleaner) TruncateTables() error {
	tables, err := getTables()
	if err != nil {
		return err
	}

	for _, table := range tables {
		if _, err = db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)); err != nil {
			return err
		}
	}

	return nil
}

func getTables() ([]string, error) {
	tables := make([]string, 0)

	rows, err := db.Query("SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public';")
	if err != nil {
		return tables, err
	}
	defer rows.Close()
	for rows.Next() {
		var value string
		if err = rows.Scan(&value); err == nil {
			tables = append(tables, value)
		}
	}

	return tables, nil
}
