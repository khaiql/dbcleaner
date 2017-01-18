package postgres

import "testing"

func TestGetTableQuery(t *testing.T) {
	helper := PostgresHelper{}
	query := helper.GetTablesQuery()

	if query != "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public';" {
		t.Error("Wrong query")
	}
}

func TestTruncateTableCommand(t *testing.T) {
	helper := PostgresHelper{}
	cmd := helper.TruncateTableCommand("users")

	if cmd != "TRUNCATE TABLE users" {
		t.Error("Wrong command")
	}
}
