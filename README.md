# DbCleaner

[![Build Status](https://travis-ci.org/khaiql/dbcleaner.svg?branch=master)](https://travis-ci.org/khaiql/dbcleaner) [![GoDoc](https://godoc.org/github.com/khaiql/dbcleaner?status.svg)](https://godoc.org/gopkg.in/khaiql/dbcleaner.v2) [![Go Report Card](https://goreportcard.com/badge/github.com/khaiql/dbcleaner)](https://goreportcard.com/report/github.com/khaiql/dbcleaner)

Clean database for testing, inspired by [database_cleaner](https://github.com/DatabaseCleaner/database_cleaner) for Ruby. It uses flock syscall under the hood to make sure the test can runs in parallel without racing issues.

## Basic usage

* To get the package, execute:

```bash
go get gopkg.in/khaiql/dbcleaner.v2
```

* To import this package, add the following line to your code:

```go
import "gopkg.in/khaiql/dbcleaner.v2"
```

* To install `TestSuite`:

```bash
go get github.com/stretchr/testify
```

* For people who is using old version (v1.0), please change your import to
```go
import "gopkg.in/khaiql/dbcleaner.v1"
```

## Options

During running test suites, there might be deadlock when 2 suites try to acquire the same table. Dbcleaner tries to
mitigate the issue by providing options for retry and panic when the deadlock couldn't be resolved after excessive retries.

```go
type Options struct {
	Logger        logging.Logger
	LockTimeout   time.Duration
	NumberOfRetry int
	RetryInterval time.Duration
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

cleaner := dbcleaner.New(SetNumberOfRetry(10), SetLockTimeout(5*time.Second))
```

## Using with testify's suite

```go
import (
	"testing"

  	"gopkg.in/khaiql/dbcleaner.v2"
  	"gopkg.in/khaiql/dbcleaner.v2/engine"
	"github.com/stretchr/testify/suite"
)

var Cleaner = dbcleaner.New()

type ExampleSuite struct {
	suite.Suite
}

func (suite *ExampleSuite) SetupSuite() {
  	// Init and set mysql cleanup engine
  	mysql := engine.NewMySQLEngine("YOUR_DB_DSN")
  	Cleaner.SetEngine(mysql)
}

func (suite *ExampleSuite) SetupTest() {
  	Cleaner.Acquire("users")
}

func (suite *ExampleSuite) TearDownTest() {
  	Cleaner.Clean("users")
}

func (suite *ExampleSuite) TestSomething() {
  	// Have some meaningful test
  	suite.Equal(true, true)
}

func TestRunSuite(t *testing.T) {
  	suite.Run(t, new(ExampleSuite))
}
```

## Support drivers

* postgres
* mysql
* sqlite3

## Write cleaner for other drivers

Basically all drivers supported by `database/sql` package are also supported by
`dbcleaner`. Check list of drivers:
[https://github.com/golang/go/wiki/SQLDrivers](https://github.com/golang/go/wiki/SQLDrivers)

For custom driver, implement your own `engine.Engine` interface and call `SetEngine` on `dbcleaner.Cleaner` instance.

## License

MIT
