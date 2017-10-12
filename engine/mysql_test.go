package engine

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMysqlTruncate(t *testing.T) {
	assert := assert.New(t)
	dsn := "root@/dbcleaner_test"
	dbEngine := NewMySQLEngine(dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}

	if err := TestInsertData(db); err != nil {
		t.Fatal(err)
	}

	t.Run("Truncate table has foreign key", func(t *testing.T) {
		err = dbEngine.Truncate("users")
		assert.NoError(err)

		count, err := TestCountRecord(db, "users")
		assert.NoError(err)
		assert.Empty(count)
	})

	t.Run("Truncate", func(t *testing.T) {
		err = dbEngine.Truncate("addresses")
		if err != nil {
			t.Error("Should truncate without error")
		}
	})
}
