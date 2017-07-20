package postgres_test

import (
	"testing"

	"github.com/khaiql/dbcleaner"
	"github.com/khaiql/dbcleaner/helper/postgres"
)

func TestGetTableQuery(t *testing.T) {
	helper := postgres.Helper{}
	query := helper.GetTablesQuery()

	if query != "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public';" {
		t.Error("Wrong query")
	}
}

func TestTruncateTableCommand(t *testing.T) {
	helper := postgres.Helper{}
	table := "users"
	cmd := helper.TruncateTableCommand(table)

	if cmd != "TRUNCATE TABLE users CASCADE" {
		t.Error("Wrong command")
	}
}

func TestInit(t *testing.T) {
	_, err := dbcleaner.FindHelper("postgres")
	if err != nil {
		t.Errorf("Shouldn't get error but got %s", err.Error())
	}
}
