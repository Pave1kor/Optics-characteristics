package models

import "database/sql"

type DBManager struct {
	Db *sql.DB
}
