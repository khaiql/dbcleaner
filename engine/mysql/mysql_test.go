package mysql

import "testing"

func TestMysqlTruncate(t *testing.T) {
	dsn := "root@/dbcleaner_test"
	engine := New(dsn)

	t.Run("Truncate table has foreign key", func(t *testing.T) {
		err := engine.Truncate("users")
		if err != nil {
			t.Error("Should truncate without error")
		}
	})

	t.Run("Truncate", func(t *testing.T) {
		err := engine.Truncate("addresses")
		if err != nil {
			t.Error("Should truncate without error")
		}
	})
}
