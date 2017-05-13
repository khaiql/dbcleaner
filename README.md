# DbCleaner

[![Build Status](https://travis-ci.org/khaiql/dbcleaner.svg?branch=master)](https://travis-ci.org/khaiql/dbcleaner) [![GoDoc](https://godoc.org/github.com/khaiql/dbcleaner?status.svg)](https://godoc.org/github.com/khaiql/dbcleaner) [![Go Report Card](https://goreportcard.com/badge/github.com/khaiql/dbcleaner)](https://goreportcard.com/report/github.com/khaiql/dbcleaner)[![Coverage Status](https://coveralls.io/repos/github/khaiql/dbcleaner/badge.svg)](https://coveralls.io/github/khaiql/dbcleaner)

Clean database for testing, inspired by [database_cleaner](https://github.com/DatabaseCleaner/database_cleaner) for Ruby

## Basic usage

* Getting started: `go get -u github.com/khaiql/dbcleaner`

```
import (
  "os"
  "testing"

  "github.com/khaiql/dbcleaner"

  // Register postgres db driver, ignore this if you have already called it
  somewhere else
  _ "github.com/lib/pq"

  // Register postgres cleaner helper
  _ "github.com/khaiql/dbcleaner/helper/postgres"

)

func TestMain(m *testing.Main) {
  cleaner, err := dbcleaner.New("postgres", "YOUR_DB_CONNECTION_STRING")
  if err != nil {
    panic(err)
  }
  defer cleaner.Close()

  code := m.Run()
  cleaner.TruncateTablesExclude("migrations")

  os.Exit(code)
}

func TestSomething(t *testing.T) {
  // TODO: Write your db related test
}
```

**NOTE:** using `TestMain` will only clear database once after all test cases

## Using with testify's suite (recommended)

```
import (
	"testing"

	"github.com/khaiql/dbcleaner"
	_ "github.com/khaiql/dbcleaner/helper/postgres"
	"github.com/stretchr/testify/suite"
)

type ExampleSuite struct {
	suite.Suite
	DBCleaner *dbcleaner.DBCleaner
}

// Init dbcleaner instance at the beginning of every suite
func (suite *ExampleSuite) SetupSuite() {
	cleaner, err := dbcleaner.New("postgres", "YOUR_DB_CONNECTION_STRING")
	if err != nil {
		panic(err)
	}

	suite.DBCleaner = cleaner
}

// Close and release connection at the end of suite
func (suite *ExampleSuite) TearDownSuite() {
	suite.DBCleaner.Close()
}

// Truncate tables after every test case. Note: sub-test using t.Run wouldn't be
// taken into account
func (suite *ExampleSuite) TearDownTest() {
	suite.DBCleaner.TruncateTablesExclude("migrations")
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

The mechanism is literally the same as `sql.RegisterDriver`. All you need is to
implement `helper.Helper` interface and call `dbcleaner.RegisterHelper`

Want example? Check [this](https://github.com/khaiql/dbcleaner/tree/master/helper/pq)

Please feel free to create PR for integrating more db drivers

## Running test

1. `docker-compose up -d`
1. `go get -u github.com/lib/pq github.com/go-sql-driver/mysql`

## License

MIT
