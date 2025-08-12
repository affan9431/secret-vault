package storage

import "database/sql"

var DB *sql.DB

func InitDB(database *sql.DB) {
	DB = database
}