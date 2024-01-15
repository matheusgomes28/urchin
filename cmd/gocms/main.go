package main

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/matheusgomes28/app"
	. "github.com/matheusgomes28/app"
	"github.com/matheusgomes28/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)
	
func loadSettings() (AppSettings, error) {
	// Want to load the environment variables
	database_address := os.Getenv("GOCMS_DATABASE_ADDRESS")
	if len(database_address) == 0 {
		return AppSettings{}, fmt.Errorf("GOCMS_DATABASE_ADDRESS is not defined\n")
	}

	database_user := os.Getenv("GOCMS_DATABASE_USER")
	if len(database_user) == 0 {
		return AppSettings{}, fmt.Errorf("GOCMS_DATABASE_USER is not defined\n")
	}

	database_password := os.Getenv("GOCMS_DATABASE_PASSWORD")
	if len(database_password) == 0 {
		return AppSettings{}, fmt.Errorf("GOCMS_DATABASE_PASSWORD is not defined\n")
	}

	database_port_str := os.Getenv("GOCMS_DATABASE_PORT")
	if len(database_port_str) == 0 {
		return AppSettings{}, fmt.Errorf("GOCMS_DATABASE_PORT is not defined\n")
	}

	database_port, err := strconv.Atoi(database_port_str)
	if err != nil {
		return AppSettings{}, fmt.Errorf("GOCMS_DATABASE_PORT is not a valid integer: %v\n", err)
	}

	return AppSettings{
		Database_user: database_user,
		Database_password: database_password,
		Database_address: database_address,
		Database_port: database_port,
	}, nil
}

func setupLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Logger created")
}

func main() {
	/// Load global application settings
	app_settings, err := loadSettings()
	if err != nil {
		log.Error().Msgf("could not get app settings: %v\n", err)
		return
	}

	db_connection, err := database.MakeSqlConnection(
		app_settings.Database_user,
		app_settings.Database_password,
		app_settings.Database_address,
		app_settings.Database_port,
	)
	if err != nil {
		log.Error().Msgf("could not create database connection: %v", err)
	}

	app.Run(app_settings, db_connection)
}
