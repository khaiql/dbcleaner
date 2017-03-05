package dbcleaner_test

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/khaiql/dbcleaner"
	"github.com/khaiql/dbcleaner/helper/mysql"
)

const (
	mysqlConnWithoutDatabase = "root:1234@/mysql"
	mysqlConnWithDatabase    = "root:1234@/dbcleaner"
	mysqlDriver              = "mysql"
)

func TestMysqlCleaner(t *testing.T) {
	dbcleaner.RegisterHelper("mysql", mysql.Helper{})
	cleaner, _ := dbcleaner.New("mysql", mysqlConnWithDatabase)
	defer cleaner.Close()

	setupMysql()
	defer dropDatabase(mysqlDriver, mysqlConnWithDatabase)

	if err := cleaner.TruncateTablesExclude("addresses"); err != nil {
		t.Errorf("Shouldn't get error, but got: %s", err.Error())
	}

	db := getDbConnection(mysqlDriver, mysqlConnWithDatabase)

	testcases := []struct {
		table         string
		expectedCount int
	}{
		{"users", 0},
		{"addresses", 1},
	}

	for _, cs := range testcases {
		count := -2
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", cs.table)
		if err := db.QueryRow(query).Scan(&count); err != nil {
			t.Fatalf("Couldn't count %s - mysql. Err: %s", cs.table, err.Error())
		}

		if count != cs.expectedCount {
			t.Errorf("Should get %d. Got %d", cs.expectedCount, count)
		}
	}
}

func setupMysql() {
	createDatabase(mysqlDriver, mysqlConnWithoutDatabase)
	db := getDbConnection(mysqlDriver, mysqlConnWithDatabase)
	defer db.Close()

	commands := []string{
		"CREATE TABLE users(name varchar(100))",
		"CREATE TABLE addresses(addr varchar(100))",
		"INSERT INTO users values(\"Dummy\")",
		"INSERT INTO addresses values(\"Singapore\")",
	}

	for _, cmd := range commands {
		_, err := db.Exec(cmd)
		if err != nil {
			panic(err)
		}
	}
}
