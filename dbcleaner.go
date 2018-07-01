// Package dbcleaner helps cleaning up database's tables upon unit test.
// With the help of https://github.com/stretchr/testify/tree/master/suite, we can easily
// acquire the tables using in the test in SetupTest or SetupSuite, and cleanup all data
// using TearDownTest or TearDownSuite
package dbcleaner

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	filemutex "github.com/alexflint/go-filemutex"
	"github.com/khaiql/dbcleaner/engine"
	"github.com/khaiql/dbcleaner/logging"
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

var (
	// ErrTableNeverLockBefore is paniced if calling Release on table that havent' been acquired before
	ErrTableNeverLockBefore = errors.New("table has never been locked before")
)

// New returns a default Cleaner with Noop Engine. Call SetEngine to set an actual working engine
func New(opts ...Option) DbCleaner {
	options := &Options{
		Logger:        &logging.Noop{},
		LockTimeout:   10 * time.Second,
		NumberOfRetry: 5,
		RetryInterval: 10 * time.Second,
	}

	for _, opt := range opts {
		opt(options)
	}

	return &cleanerImpl{
		locks:    sync.Map{},
		dbEngine: &engine.NoOp{},
		options:  options,
	}
}

type cleanerImpl struct {
	locks    sync.Map
	dbEngine engine.Engine
	options  *Options
}

func (c *cleanerImpl) loadFileMutexForTable(table string) (*filemutex.FileMutex, error) {
	fmutex, err := filemutex.New("/tmp/" + table + ".lock")
	if err != nil {
		return nil, err
	}

	value, _ := c.locks.LoadOrStore(table, fmutex)
	return value.(*filemutex.FileMutex), nil
}

func (c *cleanerImpl) SetEngine(dbEngine engine.Engine) {
	c.dbEngine = dbEngine
}

func (c *cleanerImpl) acquireTable(ctx context.Context, table string) error {
	c.options.Logger.Println("Acquiring table %s", table)

	f := func(locker *filemutex.FileMutex, d chan struct{}) {
		locker.Lock()
		d <- struct{}{}
	}

	if err := c.actOnTable(table, f); err != nil {
		return fmt.Errorf("error %s on acquire lock for table %s", err.Error(), table)
	}

	c.options.Logger.Println("Acquired lock on table %s", table)
	return nil
}

func (c *cleanerImpl) actOnTable(table string, f func(locker *filemutex.FileMutex, doneChan chan struct{})) error {
	locker, err := c.loadFileMutexForTable(table)
	if err != nil {
		panic(err)
	}
	done := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), c.options.LockTimeout)
	defer cancel()

	go f(locker, done)

	for {
		select {
		case <-done:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *cleanerImpl) releaseTable(table string) error {
	f := func(locker *filemutex.FileMutex, done chan struct{}) {
		locker.Unlock()
		locker.Close()
		c.locks.Delete(table)
		done <- struct{}{}
	}

	if err := c.actOnTable(table, f); err != nil {
		return fmt.Errorf("releaseTable %s error=%s", table, err.Error())
	}
	return nil
}

func (c *cleanerImpl) Acquire(tables ...string) {
	tried := 0

	for tried < c.options.NumberOfRetry {
		tried++
		c.options.Logger.Println("Trying to acquire %d times\n", tried)
		var err error
		acquiredTables := []string{}

		for _, table := range tables {
			err = c.acquireTable(context.Background(), table)
			if err != nil {
				break
			}
			acquiredTables = append(acquiredTables, table)
		}

		if err == nil {
			return
		}

		c.options.Logger.Println("Failed to acquired with error=%s", err.Error())

		for _, table := range acquiredTables {
			c.releaseTable(table)
		}
		time.Sleep(c.options.RetryInterval)
	}

	panic(fmt.Errorf("failed to ACQUIRE tables %v after %d times", tables, tried))
}

func (c *cleanerImpl) Clean(tables ...string) {
	for _, table := range tables {
		c.options.Logger.Println("Truncate table %s", table)
		if err := c.dbEngine.Truncate(table); err != nil {
			panic(err)
		}

		if err := c.releaseTable(table); err != nil {
			panic(err)
		}
		c.options.Logger.Println("Released lock for table %s", table)
	}
}

func (c *cleanerImpl) Close() error {
	return c.dbEngine.Close()
}
