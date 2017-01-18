package postgres

import "fmt"

type PostgresHelper struct{}

const queryAllTables = "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public';"

func (PostgresHelper) GetTablesQuery() string {
	return queryAllTables
}

func (PostgresHelper) TruncateTableCommand(tableName string) string {
	return fmt.Sprintf("TRUNCATE TABLE %s", tableName)
}
