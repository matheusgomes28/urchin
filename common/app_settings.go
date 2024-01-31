package common

import (
	"fmt"
	"os"
	"strconv"
)

type AppSettings struct {
    Database_address string
    Database_port  int
	Database_user string
	Database_password string
	Database_name string
}

func LoadSettings() (AppSettings, error) {
	// Want to load the environment variables
	database_address := os.Getenv("URCHIN_DATABASE_ADDRESS")
	if len(database_address) == 0 {
		return AppSettings{}, fmt.Errorf("URCHIN_DATABASE_ADDRESS is not defined\n")
	}

	database_user := os.Getenv("URCHIN_DATABASE_USER")
	if len(database_user) == 0 {
		return AppSettings{}, fmt.Errorf("URCHIN_DATABASE_USER is not defined\n")
	}

	database_password := os.Getenv("URCHIN_DATABASE_PASSWORD")
	if len(database_password) == 0 {
		return AppSettings{}, fmt.Errorf("URCHIN_DATABASE_PASSWORD is not defined\n")
	}

	database_name := os.Getenv("URCHIN_DATABASE_NAME")
	if len(database_name) == 0 {
		return AppSettings{}, fmt.Errorf("URCHIN_DATABASE_NAME is not defined\n")
	}

	database_port_str := os.Getenv("URCHIN_DATABASE_PORT")
	if len(database_port_str) == 0 {
		return AppSettings{}, fmt.Errorf("URCHIN_DATABASE_PORT is not defined\n")
	}
	
	database_port, err := strconv.Atoi(database_port_str)
	if err != nil {
		return AppSettings{}, fmt.Errorf("URCHIN_DATABASE_PORT is not a valid integer: %v\n", err)
	}

	return AppSettings{
		Database_user:     database_user,
		Database_password: database_password,
		Database_address:  database_address,
		Database_port:     database_port,
		Database_name: database_name,
	}, nil
}
