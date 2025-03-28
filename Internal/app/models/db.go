package models

import (
	"database/sql"
)

type DBManager struct {
	db *sql.DB
}
