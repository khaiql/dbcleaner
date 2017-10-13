package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresTruncate(t *testing.T) {
	assert := assert.New(t)
	dsn := "postgres://postgres@localhost/dbcleaner_test?sslmode=disable"
	dbEngine := NewPostgresEngine(dsn)

	t.Run("Truncate users table", func(t *testing.T) {
		err := dbEngine.Truncate("users")
		assert.NoError(err)
	})

	t.Run("Truncate addresses table", func(t *testing.T) {
		err := dbEngine.Truncate("addresses")
		assert.NoError(err)
	})
}
