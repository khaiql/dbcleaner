package engine

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestSQLiteTruncate(t *testing.T) {
	assert := assert.New(t)
	dbFilePath := "../dbcleaner_test.db"
	dbEngine := NewSqliteEngine(dbFilePath)
	db, _ := sql.Open("sqlite3", "../dbcleaner_test.db")

	t.Run("Truncate users table", func(t *testing.T) {
		err := dbEngine.Truncate("users")
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
