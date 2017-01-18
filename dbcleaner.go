package dbcleaner

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/khaiql/dbcleaner/helper"
	"github.com/khaiql/dbcleaner/utils"
)

type dbcleaner struct {
	db     *sql.DB
	driver string
}

var (
	mutex               sync.Mutex
	registeredHelpers   = make(map[string]helper.Helper)
	NotFoundHelperError = errors.New("Helper has not been registered")
)

func RegisterHelper(driverName string, helper helper.Helper) {
	mutex.Lock()
	registeredHelpers[driverName] = helper
	mutex.Unlock()
}

func New(driver, connectionString string) (*dbcleaner, error) {
	db, err := sql.Open(driver, connectionString)

	if err != nil {
		return nil, err
	}

	return &dbcleaner{db, driver}, err
}

func FindHelper(driver string) (helper.Helper, error) {
	if helper, ok := registeredHelpers[driver]; ok {
		return helper, nil
	}

	return nil, NotFoundHelperError
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
	helper, err := FindHelper(c.driver)
	if err != nil {
		return tables, err
	}

	rows, err := c.db.Query(helper.GetTablesQuery())
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
