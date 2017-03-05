package mysql_test

import (
	"testing"

	"github.com/khaiql/dbcleaner/helper/mysql"
)

func TestGetTableQuery(t *testing.T) {
	h := mysql.Helper{}
	query := h.GetTablesQuery()
	if query != "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()" {
		t.Error("Wrong query")
	}
}

func TruncateTablesCommand(t *testing.T) {
	h := mysql.Helper{}
	cmd := h.TruncateTablesCommand([]string{"users", "addresses"})
	expectedCmd := "truncate table users;truncate table addresses"

	if cmd != expectedCmd {
		t.Errorf("Expected %s. Got %s", expectedCmd, cmd)
	}
}
