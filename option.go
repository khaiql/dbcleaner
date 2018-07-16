package dbcleaner

import (
	"time"

	"github.com/khaiql/dbcleaner/logging"
)

// Options defines properties that dbEngine would use during trying to acquire and clean tables, including
// Logger: default to Noop
// LockTimeout: max duration while trying to acquire lock for a table
// NumberOfRetry: max number of retry when failed to acquire the table
// RetryInterval: sleep between each retry until reach NumberOfRetry
type Options struct {
	Logger        logging.Logger
	LockTimeout   time.Duration
	NumberOfRetry int
	RetryInterval time.Duration
	LockFileDir   string
}

type Option func(opt *Options)

// SetLogger to an instance of logging.Logger, default to Noop
func SetLogger(logger logging.Logger) Option {
	return func(opt *Options) {
		opt.Logger = logger
	}
}

// SetLockTimeout sets timeout for locking operation, default to 10 seconds
func SetLockTimeout(d time.Duration) Option {
	return func(opt *Options) {
		opt.LockTimeout = d
	}
}

// SetNumberOfRetry sets max retries for acquire the table, default to 5 times
func SetNumberOfRetry(t int) Option {
	return func(opt *Options) {
		opt.NumberOfRetry = t
	}
}

// SetRetryInterval sets sleep duration between each retry, default to 10 seconds
func SetRetryInterval(d time.Duration) Option {
	return func(opt *Options) {
		opt.RetryInterval = d
	}
}

// SetLockFileDir sets directory for lock files
func SetLockFileDir(dir string) Option {
	return func(opt *Options) {
		opt.LockFileDir = dir
	}
}
