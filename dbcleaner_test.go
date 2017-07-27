package dbcleaner_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/khaiql/dbcleaner"
)

type expectedResult struct {
	table      string
	numRecords int
}

func checkResult(db *sql.DB, expected expectedResult) error {
	numRecords, err := countRecords(db, expected.table)
	if err != nil {
		return err
	}

	if numRecords != expected.numRecords {
		return fmt.Errorf("Table %s should have %d records. Got %d", expected.table, expected.numRecords, numRecords)
	}

	return nil
}

func countRecords(db *sql.DB, table string) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	var count int
	if err := db.QueryRow(query).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func getDbConnection(driver, conn string) *sql.DB {
	db, err := sql.Open(driver, conn)
	if err != nil {
		panic(err)
	}

	return db
}

func createDatabase(driver, conn string) {
	db := getDbConnection(driver, conn)
	defer db.Close()

	_, err := db.Exec("CREATE DATABASE dbcleaner")
	if err != nil {
		panic(err)
	}
}

func dropDatabase(driver, conn string) {
	db := getDbConnection(driver, conn)
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

func (dummyHelper) TruncateTableCommand(string) string {
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
