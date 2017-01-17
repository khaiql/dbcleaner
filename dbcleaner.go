package dbcleaner

import (
	"database/sql"
	"fmt"
	"sync"

	"dbcleaner/utils"
)

type dbcleaner struct {
	db *sql.DB
}

func New(driver, connectionString string) (*dbcleaner, error) {
	db, err := sql.Open(driver, connectionString)

	if err != nil {
		return nil, err
	}

	return &dbcleaner{db}, err
}

func (c *dbcleaner) Close() error {
	return c.db.Close()
}

func (c *dbcleaner) TruncateTables(excludedTables ...string) error {
	var waitGroup sync.WaitGroup

	tables, err := c.getTables()
	if err != nil {
		return err
	}

	tables = utils.SubtractStringArray(tables, excludedTables)

	waitGroup.Add(len(tables))

	for _, table := range tables {
		go func(table string) {
			c.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
			waitGroup.Done()
		}(table)
	}

	waitGroup.Wait()

	return nil
}

func (c *dbcleaner) getTables() ([]string, error) {
	tables := make([]string, 0)

	rows, err := c.db.Query("SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public';")
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
