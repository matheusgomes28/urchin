package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/matheusgomes28/app"
	"github.com/matheusgomes28/common"
	"github.com/matheusgomes28/database"
	"github.com/rs/zerolog/log"
)

func main() {
	/// Load global application settings
	app_settings, err := common.LoadSettings()
	if err != nil {
		log.Error().Msgf("could not get app settings: %v\n", err)
		return
	}

	db_connection, err := database.MakeSqlConnection(
		app_settings.DatabaseUser,
		app_settings.DatabasePassword,
		app_settings.DatabaseAddress,
		app_settings.DatabasePort,
		app_settings.DatabaseName,
	)
	if err != nil {
		log.Error().Msgf("could not create database connection: %v", err)
	}

	app.Run(app_settings, &db_connection)
}
