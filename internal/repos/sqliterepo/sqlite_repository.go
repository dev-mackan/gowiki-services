package sqliterepo

import "database/sql"

type SqliteRepository struct {
	db *sql.DB
}
