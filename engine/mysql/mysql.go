package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	db *sql.DB
}

func New(dsn string) *Mysql {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	return &Mysql{
		db: db,
	}
}

func (mysql *Mysql) Truncate(table string) error {
	tx, err := mysql.db.Begin()
	if err != nil {
		return err
	}

	cmds := []string{
		"SET FOREIGN_KEY_CHECKS = 0",
		fmt.Sprintf("TRUNCATE %s", table),
		"SET FOREIGN_KEY_CHECKS = 1",
	}

	for _, cmd := range cmds {
		fmt.Println("Executing: ", cmd)

		if _, err := tx.Exec(cmd); err != nil {
			return tx.Rollback()
		}
	}

	return tx.Commit()
}
