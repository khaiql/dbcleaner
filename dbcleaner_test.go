package dbcleaner

import (
	"errors"
	"testing"
	"time"

	"github.com/khaiql/dbcleaner/engine"
	"github.com/stretchr/testify/mock"
)

func TestClean(t *testing.T) {
	mockEngine := &engine.MockEngine{}
	mockEngine.On("Truncate", mock.AnythingOfType("string")).Return(nil)

	Default.SetEngine(mockEngine)

	t.Run("TestNothingLock", func(t *testing.T) {
		Default.Clean("table1", "table2")
		mockEngine.AssertNumberOfCalls(t, "Truncate", 2)
		mockEngine.AssertCalled(t, "Truncate", "table1")
		mockEngine.AssertCalled(t, "Truncate", "table2")
	})

	t.Run("TestLockAndThenUnlock", func(t *testing.T) {
		tbName := "lock_table"
		Default.Acquire(tbName)
		go func() {
			time.Sleep(1 * time.Second)
			Default.Release(tbName)
		}()

		err := Default.Clean(tbName)
		if err != nil {
			t.Error(err.Error())
		}

		mockEngine.AssertCalled(t, "Truncate", tbName)
	})

	t.Run("TestTruncateError", func(t *testing.T) {
		errorTruncateMock := &engine.MockEngine{}
		errorTruncateMock.On("Truncate", mock.AnythingOfType("string")).Return(errors.New("Truncate error"))

		Default.SetEngine(errorTruncateMock)
		err := Default.Clean("error_table")

		if err.Error() != "Truncate error" {
			t.Error("Error mismatch")
		}
	})
}
