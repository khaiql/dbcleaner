# DbCleaner

[![Build Status](https://travis-ci.org/khaiql/dbcleaner.svg?branch=master)](https://travis-ci.org/khaiql/dbcleaner) [![GoDoc](https://godoc.org/github.com/khaiql/dbcleaner?status.svg)](https://godoc.org/github.com/khaiql/dbcleaner) [![Go Report Card](https://goreportcard.com/badge/github.com/khaiql/dbcleaner)](https://goreportcard.com/report/github.com/khaiql/dbcleaner)[![Coverage Status](https://coveralls.io/repos/github/khaiql/dbcleaner/badge.svg)](https://coveralls.io/github/khaiql/dbcleaner)

Clean database for testing, inspired by [database_cleaner](https://github.com/DatabaseCleaner/database_cleaner) for Ruby

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

## Using with testify's suite

```
import (
	"testing"

 "gopkg.in/khaiql/dbcleaner.v2"
	"github.com/stretchr/testify/suite"
)

type ExampleSuite struct {
	suite.Suite
}

func (suite *ExampleSuite) SetupSuite() {
  // Init and set mysql cleanup engine
  mysql := engine.NewMySQLEngine("YOUR_DB_DSN")
  dbcleanerCleaner.SetEngine(mysql)
}

func (suite *ExampleSuite) SetupTest() {
  dbcleaner.Cleaner.Acquire("users")
}

func (suite *ExampleSuite) TearDownTest() {
  dbcleaner.Cleaner.Clean("users")
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

## Write cleaner for other drivers

Basically all drivers supported by `database/sql` package are also supported by
`dbcleaner`. Check list of drivers:
[https://github.com/golang/go/wiki/SQLDrivers](https://github.com/golang/go/wiki/SQLDrivers)

For custom driver, implement your own `engine.Engine` interface and call `SetEngine` on `dbcleaner.Cleaner` instance.

## License

MIT
