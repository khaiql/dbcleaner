package engine

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPostgresTruncate(t *testing.T) {
	assert := assert.New(t)
	dsn := "postgres://postgres@localhost/dbcleaner_test?sslmode=disable"
	dbEngine := NewPostgresEngine(dsn)

	t.Run("Truncate users table", func(t *testing.T) {
		err := dbEngine.Truncate("users")
		db, _ := sql.Open("postgres", dsn)
		result, _ := db.Exec("select id from users")
		actual, _ := result.LastInsertId()

		assert.Equal(int64(0), actual)
		assert.NoError(err)
	})

	t.Run("Truncate addresses table", func(t *testing.T) {
		err := dbEngine.Truncate("addresses")
		assert.NoError(err)
	})

	t.Run("Close db connection", func(t *testing.T) {
		err := dbEngine.Close()
		assert.NoError(err)
	})
}
