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
}

type Option func(opt *Options)

func SetLogger(logger logging.Logger) Option {
	return func(opt *Options) {
		opt.Logger = logger
	}
}

func SetLockTimeout(d time.Duration) Option {
	return func(opt *Options) {
		opt.LockTimeout = d
	}
}

func SetNumberOfRetry(t int) Option {
	return func(opt *Options) {
		opt.NumberOfRetry = t
	}
}

func SetRetryInterval(d time.Duration) Option {
	return func(opt *Options) {
		opt.RetryInterval = d
	}
}
