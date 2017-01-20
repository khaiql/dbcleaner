package dbcleaner_test

import (
	"database/sql"
	"github.com/khaiql/dbcleaner"
	"github.com/khaiql/dbcleaner/helper/pq"
	_ "github.com/lib/pq"
	"testing"
)

const (
	postgresConnWithDatabaseName    = "host=127.0.0.1 port=5432 password=1234 user=postgres sslmode=disable dbname=dbcleaner"
	postgresConnWithoutDatabaseName = "host=127.0.0.1 port=5432 password=1234 user=postgres sslmode=disable"
)

func getDbConnection(conn string) *sql.DB {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	return db
}

func createDatabase(conn string) {
	db := getDbConnection(conn)
	defer db.Close()

	_, err := db.Exec("CREATE DATABASE dbcleaner")
	if err != nil {
		panic(err)
	}
}

func dropDatabase(conn string) {
	db := getDbConnection(conn)
	defer db.Close()

	_, err := db.Exec("DROP DATABASE dbcleaner")
	if err != nil {
		panic(err)
	}
}

func TestNewCleaner(t *testing.T) {
	t.Run("TestRegisteredDriver", func(t *testing.T) {
		_, err := dbcleaner.New("postgres", postgresConnWithDatabaseName)

		if err != nil {
			t.Fatalf("Should be able to open connection to db. Err: %s", err.Error())
		}
	})

	t.Run("TestUnregisteredDriver", func(t *testing.T) {
		cleaner, _ := dbcleaner.New("driver", "")

		if cleaner != nil {
			t.Fatal("Should return nil instance of cleaner")
		}
	})
}

func TestClose(t *testing.T) {
	t.Run("Connection hasn't been closed", func(t *testing.T) {
		cleaner, _ := dbcleaner.New("postgres", postgresConnWithDatabaseName)
		err := cleaner.Close()

		if err != nil {
			t.Fatalf("Should be able to close connection to database. Err: %s", err.Error())
		}
	})
}

func TestTruncateTables(t *testing.T) {
	setup()
	defer dropDatabase(postgresConnWithoutDatabaseName)

	dbcleaner.RegisterHelper("postgres", postgres.Helper{})
	cleaner, _ := dbcleaner.New("postgres", postgresConnWithDatabaseName)
	defer cleaner.Close()

	db := getDbConnection(postgresConnWithDatabaseName)
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

	t.Run("WhenExcludeAllTables", func(t *testing.T) {
		if err := cleaner.TruncateTablesExclude("users", "customers"); err == nil {
			t.Fatal("Should have error when there is no table to truncate")
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

func setup() {
	createDatabase(postgresConnWithoutDatabaseName)
	db := getDbConnection(postgresConnWithDatabaseName)
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

type dummyHelper struct{}

func (dummyHelper) GetTablesQuery() string {
	return ""
}

func (dummyHelper) TruncateTablesCommand([]string) string {
	return ""
}

func TestRegisterAndFindHelper(t *testing.T) {
	dbcleaner.RegisterHelper("dummy", dummyHelper{})

	t.Run("ExistingHelper", func(t *testing.T) {
		helper, err := dbcleaner.FindHelper("dummy")
		if err != nil {
			t.Fatalf("Shouldn't return error but got %s", err.Error())
		}

		switch helper.(type) {
		case dummyHelper:
			t.Log("OK")
		default:
			t.Error("Invalid type")
		}
	})

	t.Run("NotRegisteredHelper", func(t *testing.T) {
		_, err := dbcleaner.FindHelper("whoami")
		if err == nil {
			t.Error("It should return an error")
		}
	})
}
