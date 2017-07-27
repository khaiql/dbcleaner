package dbcleaner_test

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/khaiql/dbcleaner"
	"github.com/khaiql/dbcleaner/helper/mysql"
)

const (
	mysqlConnWithoutDatabase = "root@tcp(127.0.0.1:3306)/mysql"
	mysqlConnWithDatabase    = "root@tcp(127.0.0.1:3306)/dbcleaner"
	mysqlDriver              = "mysql"
)

func TestMysqlCleaner(t *testing.T) {
	setupMysqlDatabase()
	defer dropDatabase(mysqlDriver, mysqlConnWithDatabase)

	dbcleaner.RegisterHelper("mysql", mysql.Helper{})
	cleaner, _ := dbcleaner.New("mysql", mysqlConnWithDatabase)
	defer cleaner.Close()

	db := getDbConnection(mysqlDriver, mysqlConnWithDatabase)
	defer db.Close()

	t.Run("TruncateTablesExclude", func(t *testing.T) {
		insertMysqlTestData(db)
		defer truncateMysqlTestData(db)

		if err := cleaner.TruncateTables(); err != nil {
			t.Errorf("Shouldn't get error, but got: %s", err.Error())
		}

		expectedResults := []expectedResult{
			{table: "users", numRecords: 0},
			{table: "addresses", numRecords: 0},
		}

		for _, expectedResult := range expectedResults {
			numRecords, err := countRecords(db, expectedResult.table)
			if err != nil {
				t.Fatalf("Couldn't count %s - mysql. Err: %s", expectedResult.table, err.Error())
			}

			if numRecords != expectedResult.numRecords {
				t.Errorf("Should get %d. Got %d", expectedResult.numRecords, numRecords)
			}
		}
	})

	t.Run("TruncateTablesExclude", func(t *testing.T) {
		insertMysqlTestData(db)
		defer truncateMysqlTestData(db)

		if err := cleaner.TruncateTablesExclude("addresses"); err != nil {
			t.Fatalf("Shouldn't have error but got %s", err.Error())
		}

		expectedResults := []expectedResult{
			{table: "users", numRecords: 0},
			{table: "addresses", numRecords: 1},
		}

		for _, expected := range expectedResults {
			if err := checkResult(db, expected); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("TruncateTablesOnly", func(t *testing.T) {
		insertMysqlTestData(db)
		defer truncateMysqlTestData(db)

		if err := cleaner.TruncateTablesOnly("addresses"); err != nil {
			t.Errorf("Shouldn't get error, but got: %s", err.Error())
		}

		expectedResults := []expectedResult{
			{table: "users", numRecords: 1},
			{table: "addresses", numRecords: 0},
		}

		for _, expected := range expectedResults {
			if err := checkResult(db, expected); err != nil {
				t.Error(err)
			}
		}
	})
}

func setupMysqlDatabase() {
	createDatabase(mysqlDriver, mysqlConnWithoutDatabase)
	db := getDbConnection(mysqlDriver, mysqlConnWithDatabase)
	defer db.Close()

	commands := []string{
		"CREATE TABLE users(name varchar(100))",
		"CREATE TABLE addresses(addr varchar(100))",
	}

	for _, cmd := range commands {
		_, err := db.Exec(cmd)
		if err != nil {
			panic(err)
		}
	}
}

func insertMysqlTestData(db *sql.DB) {
	commands := []string{
		"INSERT INTO addresses values(\"Singapore\")",
		"INSERT INTO users values(\"Dummy\")",
	}

	for _, cmd := range commands {
		_, err := db.Exec(cmd)
		if err != nil {
			panic(err)
		}
	}
}

func truncateMysqlTestData(db *sql.DB) {
	commands := []string{
		"TRUNCATE TABLE addresses",
		"TRUNCATE TABLE users",
	}

	for _, cmd := range commands {
		_, err := db.Exec(cmd)
		if err != nil {
			panic(err)
		}
	}
}
