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

	// Release releases locks on tables which have been acquired before
	Release(tables ...string)

	// Clean calls Truncate the tables
	Clean(tables ...string) error

	// Close calls corresponding method on dbEngine to release connection to db
	Close() error
}

// Cleaner implementation of DbCleaner. Its default dbEngine is NoOp
// Use SetEngine to set actual dbEngine that your app is using
var Cleaner DbCleaner

// ErrTableNeverLockBefore is paniced if calling Release on table that havent' been acquired before
var ErrTableNeverLockBefore = errors.New("Table has never been locked before")

func init() {
	Cleaner = &cleanerImpl{
		locks:    make(map[string]*sync.RWMutex),
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

func (c *cleanerImpl) Release(tables ...string) {
	for _, table := range tables {
		if c.locks[table] == nil {
			panic(ErrTableNeverLockBefore)
		}

		c.locks[table].RUnlock()
	}
}

func (c *cleanerImpl) Clean(tables ...string) error {
	for _, table := range tables {
		if c.locks[table] != nil {
			c.locks[table].Lock()
			defer c.locks[table].Unlock()
		}

		if err := c.dbEngine.Truncate(table); err != nil {
			return err
		}
	}

	return nil
}

func (c *cleanerImpl) Close() error {
	return c.dbEngine.Close()
}
