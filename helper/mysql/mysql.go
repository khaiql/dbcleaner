package mysql

import "fmt"
import "strings"

// Helper is a mysql helper
type Helper struct{}

const getTableQuery = "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()"

// GetTablesQuery returns a query to get all tables of connected mysql database
func (Helper) GetTablesQuery() string {
	return getTableQuery
}

// TruncateTablesCommand returns command to truncate all tables supplied in
// argument
func (Helper) TruncateTablesCommand(tableNames []string) string {
	cmds := make([]string, len(tableNames), len(tableNames))
	for _, table := range tableNames {
		cmds = append(cmds, fmt.Sprintf("TRUNCATE TABLE %s", table))
	}

	return strings.Join(cmds, ";")
}
