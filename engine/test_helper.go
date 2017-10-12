package engine

import (
	"database/sql"
	"fmt"
)

func TestInsertUser(db *sql.DB) (int64, error) {
	result, err := db.Exec("INSERT INTO users values ()")

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func TestInsertAddress(db *sql.DB, userID int64) (int64, error) {
	result, err := db.Exec(fmt.Sprintf("INSERT INTO addresses (user_id) values (%d)", userID))

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func TestInsertData(db *sql.DB) error {
	userID, err := TestInsertUser(db)
	if err != nil {
		return err
	}

	_, err = TestInsertAddress(db, userID)

	if err != nil {
		return err
	}

	return nil
}

func TestCountRecord(db *sql.DB, table string) (int, error) {
	var count int

	result := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table))

	err := result.Scan(&count)
	return count, err
}
