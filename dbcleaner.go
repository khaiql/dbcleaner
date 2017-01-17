package dbcleaner

import "database/sql"

var db *sql.DB

type dbcleaner struct{}

func New(driver, connectionString string) (*dbcleaner, error) {
	var err error
	db, err = sql.Open(driver, connectionString)

	if err != nil {
		return nil, err
	}

	return &dbcleaner{}, err
}

func (c *dbcleaner) Close() error {
	return db.Close()
}
