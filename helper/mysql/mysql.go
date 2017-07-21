package mysql

import (
	"fmt"

	"github.com/khaiql/dbcleaner"
)

// Helper is a mysql helper
type Helper struct{}

func init() {
	dbcleaner.RegisterHelper("mysql", Helper{})
}

const getTableQuery = "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()"

// GetTablesQuery returns a query to get all tables of connected mysql database
func (Helper) GetTablesQuery() string {
	return getTableQuery
}

// TruncateTableCommand returns mysql command to truncate a table
func (Helper) TruncateTableCommand(table string) string {
	return fmt.Sprintf("TRUNCATE TABLE %s;", table)
}
