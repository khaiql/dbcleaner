package postgres

import "fmt"

type Helper struct{}

const queryAllTables = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public';"

func (Helper) GetTablesQuery() string {
	return queryAllTables
}

func (Helper) TruncateTableCommand(tableName string) string {
	return fmt.Sprintf("TRUNCATE TABLE %s", tableName)
}
