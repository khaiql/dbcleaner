package dbcleaner_test

import (
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

func TestPostgresTruncateTable(t *testing.T) {
	setupPostgres()
	defer dropDatabase(postgresDriver, postgresConnWithoutDatabaseName)

	dbcleaner.RegisterHelper("postgres", postgres.Helper{})
	cleaner, _ := dbcleaner.New("postgres", postgresConnWithDatabaseName)
	defer cleaner.Close()

	db := getDbConnection(postgresDriver, postgresConnWithDatabaseName)
	defer db.Close()

	t.Run("WithoutExcludedTables", func(t *testing.T) {
		if err := cleaner.TruncateTables(); err != nil {
			t.Fatalf("Shouldn't have error but got %s", err.Error())
		}

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
		if err != nil {
			t.Fatalf("Shouldn't have error but got: %s", err.Error())
		}

		if count != 0 {
			t.Errorf("Should get 0, but got: %d", count)
		}
	})

	t.Run("WithExludedTables", func(t *testing.T) {
		db.Exec("INSERT INTO users(name) values('username')")

		if err := cleaner.TruncateTablesExclude("users"); err != nil {
			t.Fatalf("Shouldn't have error but got %s", err.Error())
		}

		var count int
		db.QueryRow("SELECT COUNT(*) FROM users;").Scan(&count)
		if count != 1 {
			t.Errorf("Should get 1 but got %d", count)
		}
	})
}

func setupPostgres() {
	createDatabase(postgresDriver, postgresConnWithoutDatabaseName)
	db := getDbConnection(postgresDriver, postgresConnWithDatabaseName)
	defer db.Close()

	commands := []string{
		"CREATE TABLE users(id serial primary key, name varchar)",
		"CREATE TABLE customers(id serial primary key, name varchar)",
		"INSERT INTO users(name) values ('UserA')",
	}

	for _, cmd := range commands {
		_, err := db.Exec(cmd)
		if err != nil {
			panic(err)
		}
	}
}
