package database

import (
	"database/sql"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func MakeSqliteConnection(databaseFile string) (SqlDatabase, error) {

	/// TODO : let user specify the DB
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		return SqlDatabase{}, err
	}

	if err := db.Ping(); err != nil {
		return SqlDatabase{}, err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Second * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return SqlDatabase{
		Address:    "",
		Port:       0,
		User:       "",
		Connection: db,
	}, nil
}
