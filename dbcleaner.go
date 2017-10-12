package dbcleaner

import (
	"sync"

	"github.com/khaiql/dbcleaner/engine"
)

type Cleaner interface {
	SetEngine(dbEngine engine.Engine)
	RLock(tables ...string)
	RUnlock(tables ...string)
	Clean(tables ...string) error
}

var DefaultCleaner Cleaner

func init() {
	DefaultCleaner = &cleanerImpl{
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

func (c *cleanerImpl) RLock(tables ...string) {
	for _, table := range tables {
		if c.locks[table] == nil {
			c.locks[table] = new(sync.RWMutex)
		}

		c.locks[table].RLock()
	}
}

func (c *cleanerImpl) RUnlock(tables ...string) {
	for _, table := range tables {
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
