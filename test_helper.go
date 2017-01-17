package dbcleaner_test

import "database/sql"

const (
	withDbNameConnection    = "host=127.0.0.1 port=5432 password=1234 user=postgres sslmode=disable dbname=dbcleaner"
	withoutDbNameConnection = "host=127.0.0.1 port=5432 password=1234 user=postgres sslmode=disable"
)

func getDbConnection(conn string) *sql.DB {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	return db
}

func createDatabase() {
	db := getDbConnection(withoutDbNameConnection)
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE dbcleaner")
	if err != nil {
		panic(err)
	}
}

func dropDatabase() {
	db := getDbConnection(withoutDbNameConnection)
	defer db.Close()

	_, err = db.Exec("DROP DATABASE dbcleaner")
	if err != nil {
		panic(err)
	}
}
