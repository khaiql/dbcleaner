package dbcleaner_test

import (
	"dbcleaner"
	"testing"

	_ "github.com/lib/pq"
)

func connectionString() string {
	return "host=localhost port=5433 password=1234 username=postgres sslmode=disabled"
}

func TestNewCleaner(t *testing.T) {
	t.Run("TestRegisteredDriver", func(t *testing.T) {
		_, err := dbcleaner.New("postgres", connectionString())

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
