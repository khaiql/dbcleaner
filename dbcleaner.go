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

func (c *cleanerImpl) loadFileMutexForTable(table string) (*filemutex.FileMutex, error) {
	value, ok := c.locks.Load(table)
	if !ok {
		fmutex, err := filemutex.New("/tmp/" + table + ".lock")
		if err == nil {
			c.locks.Store(table, fmutex)
		}

		return fmutex, err
	}

	return value.(*filemutex.FileMutex), nil
}

func (c *cleanerImpl) SetEngine(dbEngine engine.Engine) {
	c.dbEngine = dbEngine
}

func (c *cleanerImpl) Acquire(tables ...string) {
	for _, table := range tables {
		locker, err := c.loadFileMutexForTable(table)
		if err != nil {
			panic(err)
		}

		locker.Lock()
	}
}

func (c *cleanerImpl) Clean(tables ...string) {
	for _, table := range tables {
		locker, err := c.loadFileMutexForTable(table)
		if err != nil {
			panic(err)
		}

		doneChan := make(chan bool)

		go func() {
			select {
			case <-doneChan:
				return
			case <-time.After(10 * time.Second): // Timeout if couldn't acquire the lock after some time
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
	// Close all file descriptor
	c.locks.Range(func(_, value interface{}) bool {
		locker := value.(*filemutex.FileMutex)
		locker.Close()
		return true
	})

	return c.dbEngine.Close()
}
