package dbcleaner

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/khaiql/dbcleaner/helper"
	"github.com/khaiql/dbcleaner/utils"
)

type dbcleaner struct {
	db     *sql.DB
	driver string
}

var (
	mutex             sync.Mutex
	registeredHelpers = make(map[string]helper.Helper)

	// ErrHelperNotFound return when calling an unregistered Helper
	ErrHelperNotFound = errors.New("Helper has not been registered")
)

// RegisterHelper register an Helper instance for a particular driver
func RegisterHelper(driverName string, helper helper.Helper) {
	mutex.Lock()
	defer mutex.Unlock()
	registeredHelpers[driverName] = helper
}

// New returns a Cleaner instance for a particular driver
func New(driver, connectionString string) (*dbcleaner, error) {
	db, err := sql.Open(driver, connectionString)

	if err != nil {
		return nil, err
	}

	return &dbcleaner{db, driver}, err
}

// FindHelper return a registered Helper using driver name
func FindHelper(driver string) (helper.Helper, error) {
	if helper, ok := registeredHelpers[driver]; ok {
		return helper, nil
	}

	return nil, ErrHelperNotFound
}

func (c *dbcleaner) Close() error {
	return c.db.Close()
}

func (c *dbcleaner) TruncateTables() error {
	return c.TruncateTablesExclude()
}

func (c *dbcleaner) TruncateTablesExclude(excludedTables ...string) error {
	tables, err := c.getTables()
	if err != nil {
		return err
	}

	helper, err := FindHelper(c.driver)
	if err != nil {
		return err
	}

	tables = utils.SubtractStringArray(tables, excludedTables)
	_, err = c.db.Exec(helper.TruncateTablesCommand(tables))
	return err
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
