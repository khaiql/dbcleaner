package postgres

import (
	"fmt"

	"github.com/khaiql/dbcleaner"
)

// Helper postgres specific helper
type Helper struct{}

func init() {
	dbcleaner.RegisterHelper("postgres", Helper{})
}

const queryAllTables = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public';"

// GetTablesQuery returns a query string to get list of tables that should be
// truncated
func (Helper) GetTablesQuery() string {
	return queryAllTables
}

// TruncateTableCommand returns postgres command to truncate table
func (Helper) TruncateTableCommand(table string) string {
	return fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
}
