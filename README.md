# DbCleaner

[![Build Status](https://travis-ci.org/khaiql/dbcleaner.svg?branch=master)](https://travis-ci.org/khaiql/dbcleaner)

Clean database for testing, inspired by [database_cleaner](https://github.com/DatabaseCleaner/database_cleaner) for Ruby

## How to use

* Getting started: `go get -u github.com/khaiql/dbcleaner`

```
import (
  "os"
  "testing"

  "github.com/khaiql/dbcleaner"

  // Register postgres db driver
  _ "github.com/lib/pq"

  // Register postgres cleaner helper
  _ "github.com/khaiql/dbcleaner/helper/pq"

)

func TestMain(m *testing.Main) {
  cleaner, err := dbcleaner.New("postgres", "YOUR_DB_CONNECTION_STRING")
  if err != nil {
    panic(err)
  }
  defer cleaner.Close()

  exitCode = m.Run()
  // Truncate all but exclude migrations table
  cleaner.TruncateTables("migrations")
  os.Exit(exitCode)
}

func TestSomething(t *testing.T) {
  // TODO: Write your db related test
}
```

## Supporting drivers

* postgres

## Write cleaner for other drivers

Basically all drivers supported by `database/sql` package are supported by
`dbcleaner`. Find all drivers at
[https://github.com/golang/go/wiki/SQLDrivers](https://github.com/golang/go/wiki/SQLDrivers)

The mechanism is literally the same as `sql.RegisterDriver`. All you need is to
implement `helper.Helper` interface of `dbcleaner`.

Want example? Check [this](https://github.com/khaiql/dbcleaner/tree/master/helper/pq)

Please feel free to create PR for integrating more db drivers

## License

MIT
