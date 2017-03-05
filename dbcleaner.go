package dbcleaner

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/khaiql/dbcleaner/helper"
	"github.com/khaiql/dbcleaner/utils"
)

// DBCleaner instance of cleaner that can perform cleaning tables data
type DBCleaner struct {
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

// New returns a Cleaner instance for a particular driver using provided
// connectionString
func New(driver, connectionString string) (*DBCleaner, error) {
	db, err := sql.Open(driver, connectionString)

	if err != nil {
		return nil, err
	}

	return &DBCleaner{db, driver}, err
}

// FindHelper return a registered Helper using driver name
func FindHelper(driver string) (helper.Helper, error) {
	if helper, ok := registeredHelpers[driver]; ok {
		return helper, nil
	}

	return nil, ErrHelperNotFound
}

// Close closes connection to database
func (c *DBCleaner) Close() error {
	return c.db.Close()
}

// TruncateTables truncates data of all tables
func (c *DBCleaner) TruncateTables() error {
	return c.TruncateTablesExclude()
}

// TruncateTablesExclude truncates data of all tables but exclude some specify
// in the list
func (c *DBCleaner) TruncateTablesExclude(excludedTables ...string) error {
	tables, err := c.getTables()
	if err != nil {
		return err
	}

	helper, err := FindHelper(c.driver)
	if err != nil {
		return err
	}

	tables = utils.SubtractStringArray(tables, excludedTables)

	var wg sync.WaitGroup
	wg.Add(len(tables))
	for _, table := range tables {
		go func(tbl string) {
			cmd := helper.TruncateTableCommand(tbl)
			c.db.Exec(cmd)
			wg.Done()
		}(table)
	}

	return err
}

func (c *DBCleaner) getTables() ([]string, error) {
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
