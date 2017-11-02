// Package dbcleaner helps cleaning up database's tables upon unit test.
// With the help of https://github.com/stretchr/testify/tree/master/suite, we can easily
// acquire the tables using in the test in SetupTest or SetupSuite, and cleanup all data
// using TearDownTest or TearDownSuite
package dbcleaner

import (
	"errors"
	"sync"

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
		locks:    map[string]*sync.RWMutex{},
		dbEngine: &engine.NoOp{},
	}
}

type cleanerImpl struct {
	locks    map[string]*sync.RWMutex
	dbEngine engine.Engine
}

func (c *cleanerImpl) SetEngine(dbEngine engine.Engine) {
	c.dbEngine = dbEngine
}

func (c *cleanerImpl) Acquire(tables ...string) {
	for _, table := range tables {
		if c.locks[table] == nil {
			c.locks[table] = new(sync.RWMutex)
		}

		c.locks[table].RLock()
	}
}

func (c *cleanerImpl) Clean(tables ...string) {
	for _, table := range tables {
		if c.locks[table] != nil {
			c.locks[table].RUnlock()
			c.locks[table].Lock()
			defer c.locks[table].Unlock()
		}

		if err := c.dbEngine.Truncate(table); err != nil {
			panic(err)
		}
	}
}

func (c *cleanerImpl) Close() error {
	return c.dbEngine.Close()
}
