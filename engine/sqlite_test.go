package engine

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestSQLiteTruncate(t *testing.T) {
	assert := assert.New(t)
	dbFilePath := "../test.db"
	dbEngine := NewSqliteEngine(dbFilePath)

	t.Run("Truncate users table", func(t *testing.T) {
		err := dbEngine.Truncate("users")
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
