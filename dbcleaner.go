// Package dbcleaner helps cleaning up database's tables upon unit test.
// With the help of https://github.com/stretchr/testify/tree/master/suite, we can easily
// acquire the tables using in the test in SetupTest or SetupSuite, and cleanup all data
// using TearDownTest or TearDownSuite
package dbcleaner

import (
	"errors"
	"fmt"
	"sync"
	"time"

	filemutex "github.com/alexflint/go-filemutex"
	"github.com/khaiql/dbcleaner/engine"
)

// DbCleaner interface
type DbCleaner interface {
	// SetEngine sets dbEngine, can be mysql, postgres...
	SetEngine(dbEngine engine.Engine)

	// Acquire will lock tables passed in params so data in the table would not be deleted by other test cases
	Acquire(tables ...string)

	// Clean calls Truncate the tables
	Clean(tables ...string)

	// Close calls corresponding method on dbEngine to release connection to db
	Close() error
}

// ErrTableNeverLockBefore is paniced if calling Release on table that havent' been acquired before
var ErrTableNeverLockBefore = errors.New("Table has never been locked before")

// New returns a default Cleaner with Noop Engine. Call SetEngine to set an actual working engine
func New() DbCleaner {
	return &cleanerImpl{
		locks:    sync.Map{},
		dbEngine: &engine.NoOp{},
	}
}

type cleanerImpl struct {
	locks    sync.Map
	dbEngine engine.Engine
}

func (c *cleanerImpl) SetEngine(dbEngine engine.Engine) {
	c.dbEngine = dbEngine
}

func (c *cleanerImpl) Acquire(tables ...string) {
	for _, table := range tables {
		var locker *filemutex.FileMutex
		var err error

		if l, ok := c.locks.Load(table); !ok {
			locker, err = filemutex.New("/tmp/" + table + ".lock")
			if err != nil {
				panic(err)
			}

			c.locks.Store(table, locker)
		} else {
			locker = l.(*filemutex.FileMutex)
		}

		locker.Lock()
	}
}

func (c *cleanerImpl) Clean(tables ...string) {
	for _, table := range tables {
		var locker *filemutex.FileMutex
		var err error

		if l, ok := c.locks.Load(table); !ok {
			locker, err = filemutex.New("/tmp/" + table + ".lock")
			if err != nil {
				panic(err)
			}

			c.locks.Store(table, locker)
		} else {
			locker = l.(*filemutex.FileMutex)
		}

		doneChan := make(chan bool)

		go func() {
			select {
			case <-doneChan:
				return
			case <-time.After(10 * time.Second):
				panic(fmt.Sprintf("couldn't acquire the lock for table %s because of timeout", table))
			}
		}()

		if err := c.dbEngine.Truncate(table); err != nil {
			panic(err)
		}

		doneChan <- true
		locker.Unlock()
	}
}

func (c *cleanerImpl) Close() error {
	c.locks.Range(func(_, value interface{}) bool {
		locker := value.(*filemutex.FileMutex)
		locker.Close()
		return true
	})
	return c.dbEngine.Close()
}
