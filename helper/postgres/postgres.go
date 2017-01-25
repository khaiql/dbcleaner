package postgres

import (
	"fmt"
	"strings"

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

// TruncateTablesCommand accepts list of tables as argument, returns a command
// to clear all data of tables
func (Helper) TruncateTablesCommand(tableNames []string) string {
	return fmt.Sprintf("TRUNCATE TABLE %s", strings.Join(tableNames, ","))
}
