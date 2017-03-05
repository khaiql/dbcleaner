package dbcleaner_test

import (
	"database/sql"
	"testing"

	"github.com/khaiql/dbcleaner"
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
