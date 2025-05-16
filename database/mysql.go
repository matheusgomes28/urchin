package database

import (
	"database/sql"
	"fmt"
	"time"
)

func MakeMysqlConnection(user string, password string, address string, port int, database string) (SqlDatabase, error) {

	/// TODO : let user specify the DB
	connection_str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, address, port, database)
	db, err := sql.Open("mysql", connection_str)
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
		Address:    address,
		Port:       port,
		User:       user,
		Connection: db,
	}, nil
}
