package dbcleaner_test

import (
	"database/sql"
	"testing"

	"github.com/khaiql/dbcleaner"
	"github.com/khaiql/dbcleaner/helper/postgres"

	_ "github.com/lib/pq"
)

const (
	postgresConnWithDatabaseName    = "host=127.0.0.1 port=5432 password=1234 user=postgres sslmode=disable dbname=dbcleaner"
	postgresConnWithoutDatabaseName = "host=127.0.0.1 port=5432 password=1234 user=postgres sslmode=disable"
	postgresDriver                  = "postgres"
)

func TestPostgresCleaner(t *testing.T) {
	setupPostgresDatabase()
	defer dropDatabase(postgresDriver, postgresConnWithoutDatabaseName)

	dbcleaner.RegisterHelper("postgres", postgres.Helper{})
	cleaner, _ := dbcleaner.New("postgres", postgresConnWithDatabaseName)
	defer cleaner.Close()

	db := getDbConnection(postgresDriver, postgresConnWithDatabaseName)
	defer db.Close()

	t.Run("TruncateTables", func(t *testing.T) {
		insertPostgresTestData(db)
		defer truncatePostgresTestData(db)

		if err := cleaner.TruncateTables(); err != nil {
			t.Fatalf("Shouldn't have error but got %s", err.Error())
		}

		expectedResults := []expectedResult{
			{table: "users", numRecords: 0},
			{table: "customers", numRecords: 0},
		}

		for _, expected := range expectedResults {
			if err := checkResult(db, expected); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("TruncateTablesExclude", func(t *testing.T) {
		insertPostgresTestData(db)
		defer truncatePostgresTestData(db)

		if err := cleaner.TruncateTablesExclude("users"); err != nil {
			t.Fatalf("Shouldn't have error but got %s", err.Error())
		}

		expectedResults := []expectedResult{
			{table: "users", numRecords: 1},
			{table: "customers", numRecords: 0},
		}

		for _, expected := range expectedResults {
			if err := checkResult(db, expected); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("TruncateTablesOnly", func(t *testing.T) {
		insertPostgresTestData(db)
		defer truncatePostgresTestData(db)

		if err := cleaner.TruncateTablesOnly("customers"); err != nil {
			t.Fatalf("Shouldn't have error but got %s", err.Error())
		}

		expectedResults := []expectedResult{
			{table: "users", numRecords: 1},
			{table: "customers", numRecords: 0},
		}

		for _, expected := range expectedResults {
			if err := checkResult(db, expected); err != nil {
				t.Error(err)
			}
		}
	})
}

func setupPostgresDatabase() {
	createDatabase(postgresDriver, postgresConnWithoutDatabaseName)
	db := getDbConnection(postgresDriver, postgresConnWithDatabaseName)
	defer db.Close()

	commands := []string{
		"CREATE TABLE users(id serial primary key, name varchar)",
		"CREATE TABLE customers(id serial primary key, user_id integer REFERENCES users)",
	}

	for _, cmd := range commands {
		_, err := db.Exec(cmd)
		if err != nil {
			panic(err)
		}
	}
}

func insertPostgresTestData(db *sql.DB) {
	commands := []string{
		"INSERT INTO users(name) values ('UserA')",
		"INSERT INTO customers(user_id) values (currval('users_id_seq'))",
	}

	for _, cmd := range commands {
		_, err := db.Exec(cmd)
		if err != nil {
			panic(err)
		}
	}
}

func truncatePostgresTestData(db *sql.DB) {
	commands := []string{
		"TRUNCATE TABLE users CASCADE",
	}

	for _, cmd := range commands {
		_, err := db.Exec(cmd)
		if err != nil {
			panic(err)
		}
	}
}
