package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMysqlTruncate(t *testing.T) {
	assert := assert.New(t)
	dsn := "root@/dbcleaner_test"
	dbEngine := NewMySQLEngine(dsn)

	t.Run("Truncate table has foreign key", func(t *testing.T) {
		err := dbEngine.Truncate("users")
		assert.NoError(err)
	})

	t.Run("Truncate", func(t *testing.T) {
		err := dbEngine.Truncate("addresses")
		assert.NoError(err)
	})
}
