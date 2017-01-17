package dbcleaner_test

import (
	"database/sql"
	"dbcleaner"
	"testing"

	_ "github.com/lib/pq"
)

const (
	connWithDatabaseName    = "host=127.0.0.1 port=5432 password=1234 user=postgres sslmode=disable dbname=dbcleaner"
	connWithoutDatabaseName = "host=127.0.0.1 port=5432 password=1234 user=postgres sslmode=disable"
)

func getDbConnection(conn string) *sql.DB {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	return db
}

func createDatabase() {
	db := getDbConnection(connWithoutDatabaseName)
	defer db.Close()

	_, err := db.Exec("CREATE DATABASE dbcleaner")
	if err != nil {
		panic(err)
	}
}

func dropDatabase() {
	db := getDbConnection(connWithoutDatabaseName)
	defer db.Close()

	_, err := db.Exec("DROP DATABASE dbcleaner")
	if err != nil {
		panic(err)
	}
}

func TestNewCleaner(t *testing.T) {
	t.Run("TestRegisteredDriver", func(t *testing.T) {
		_, err := dbcleaner.New("postgres", connWithDatabaseName)

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
		cleaner, _ := dbcleaner.New("postgres", connWithDatabaseName)
		err := cleaner.Close()

		if err != nil {
			t.Fatalf("Should be able to close connection to database. Err: %s", err.Error())
		}
	})
}

func TestTruncateTables(t *testing.T) {
	setup()
	defer dropDatabase()

	cleaner, _ := dbcleaner.New("postgres", connWithDatabaseName)
	_ = cleaner.TruncateTables()

	db := getDbConnection(connWithDatabaseName)
	defer db.Close()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Fatalf("Shouldn't have error but got: %s", err.Error())
	}

	if count != 0 {
		t.Errorf("Should get 0, but got: %d", count)
	}
}

func setup() {
	createDatabase()
	db := getDbConnection(connWithDatabaseName)
	defer db.Close()

	commands := []string{
		"CREATE TABLE users(id serial primary key, name varchar)",
		"INSERT INTO users(name) values ('UserA')",
	}

	for _, cmd := range commands {
		_, err := db.Exec(cmd)
		if err != nil {
			panic(err)
		}
	}
}
