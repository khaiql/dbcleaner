package dbcleaner_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlConnWithoutDatabase = "root:1234@/mysql"
	mysqlConnWithDatabase    = "root:1234@/dbcleaner"
	mysqlDriver              = "mysql"
)

func TestMysqlCleaner(t *testing.T) {
	setupMysql()
}

func setupMysql() {
	createDatabase(mysqlDriver, mysqlConnWithoutDatabase)
	db := getDbConnection(mysqlDriver, mysqlConnWithDatabase)
	defer db.Close()
}
